package mydatastore

import (
	"dataaccess/models"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

const (
	MAX_CONN int           = 10
	MAX_IDLE int           = 5
	MAX_LIFE time.Duration = 0 // forever
)

type DBConfig struct {
	Uid     string
	Pwd     string
	IP      string
	Port    string
	Dbname  string
	maxConn int
	maxIdle int
	maxLife time.Duration
}

type MyDatastore struct {
	db *sql.DB
}

var mainconn *sql.DB

func (o DBConfig) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", o.Uid, o.Pwd, o.IP, o.Port, o.Dbname)
}

func (o DBConfig) checkDefaults() {
	if o.maxConn == 0 {
		o.maxConn = MAX_CONN
	}
	if o.maxIdle == 0 {
		o.maxIdle = MAX_IDLE
	}
}

func NewDatastore(o DBConfig) (*MyDatastore, error) {
	newObj := new(MyDatastore)
	var err error
	newObj.db, err = o.connect()
	if err != nil {
		return nil, err
	}
	return newObj, nil
}

func (o DBConfig) connect() (*sql.DB, error) {
	db, err := sql.Open("mysql", o.String())
	if err != nil {
		return nil, err
	}
	//defer db.Close()
	o.checkDefaults()
	db.SetMaxIdleConns(o.maxIdle)    // important when db is PaaS, to be close to 0
	db.SetConnMaxLifetime(o.maxLife) // important when db is PaaS
	db.SetMaxOpenConns(o.maxConn)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Connect and check the server version
	var version, id string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	db.QueryRow("SELECT * from table_").Scan(&id)
	fmt.Println("Table id:", id)

	mainconn = db

	return db, nil
}

// BEGIN Store functions

// StoreTable adds 1 row as
// 1. insert into table_
func (o *MyDatastore) StoreTable(t *models.Table) error {
	var sqlstr string
	var err error
	sqlstr, err = t.GetInsertTable()
	if err != nil {
		return err
	}

	_, err = o.db.Exec(sqlstr)
	if err != nil {
		return err
	}

	return nil
}

// StoreTableColnames in a Transaction:
// 1. create table owner_name_colnames
// 2. insert into owner_name_colnames
// 3. update NCols in table_
func (o *MyDatastore) StoreTableColnames(t *models.TableColnames) error {
	var sqlstr [3]string
	var err error

	sqlstr[0], err = t.GetCreateTable()
	if err != nil {
		return err
	}
	sqlstr[1], err = t.GetInsertTable()
	if err != nil {
		return err
	}
	sqlstr[2], err = t.Parent().GetUpdateNCols()
	if err != nil {
		return err
	}

	// Transaction starts
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	for _, stmt := range sqlstr {

		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()
}

// StoreTableValues in a Transaction:
// 1. create table owner_name_values
// 2. insert into owner_name_values
// 3. TODO: GetIncrementTable(affectedRows)
func (o *MyDatastore) StoreTableValues(t *models.TableValues) error {
	var sqlstr [2]string
	var err error

	sqlstr[0], err = t.GetCreateTable()
	if err != nil {
		return err
	}
	sqlstr[1], err = t.GetInsertTable()
	if err != nil {
		return err
	}

	// Transaction starts
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	for _, stmt := range sqlstr {

		_, err = tx.Exec(stmt)
		if err != nil {
			tx.Rollback()
			return err
		}

	}

	return tx.Commit()
}

// UpdateTable just updates 2 fields of table_: descr, tags
func (o *MyDatastore) UpdateTable(t *models.Table) error {
	sql, err := t.GetUpdateTable()
	if err != nil {
		return err
	}
	_, err = o.db.Exec(sql)

	return err
}

// AddColnames does in a Transaction:
// if lang exists
// 		1. delete owner_name_colnames where lang=
// 1. insert owner_name_colnames (lang, ...)
func (o *MyDatastore) AddColnames(t *models.Table) error {
	var sql1, sql2 string
	var err error
	sql1, err = t.Colnames.GetDeleteTable()
	if err != nil {
		return err
	}
	sql2, err = t.Colnames.GetInsertTable()
	if err != nil {
		return err
	}
	// Transaction starts
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(sql1)
	if err != nil {
		tx.Rollback()
		return err
	}

	_, err = tx.Exec(sql2)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

// AddValues
// 1. insert owner_name_values
// 2. update table_ (nrows++)
func (o *MyDatastore) AddValues(t *models.Table) error {
	var sql1, sql2 string
	var err error
	sql1, err = t.Values.GetInsertTable()
	if err != nil {
		return err
	}

	// Transaction starts
	// tx, errTx := mainconn.db.Begin()
	// if errTx != nil {
	// 	return err
	// }

	res, err1 := o.db.Exec(sql1)
	if err1 != nil {
		return err1
	}
	nrows, _ := res.RowsAffected()

	sql2, err = t.GetIncrementTable(nrows)
	if err != nil {
		return err
	}

	_, err = o.db.Exec(sql2)
	if err != nil {
		return err
	}

	return nil
	//return tx.Commit()
}

// END Strore functions

// START Read() functions

// ReadTable returns the models.Table without colnames neither values
func (o *MyDatastore) ReadTable(name string, owner string) (*models.Table, error) {
	sqlstr, errParam := models.GetSelectTable(name, owner)

	if errParam != nil {
		return nil, errParam
	}
	rows, err := o.db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	table := models.Table{}
	for rows.Next() {
		if err = rows.Scan(&table.Id, &table.Descr, &table.Tags, &table.DefLang, &table.NCols, &table.NRows); err != nil {
			return nil, err
		}
	}
	table.Name = name
	table.Owner = owner

	// tableColnames, errColnames := o.ReadTableColnames(table.DefLang)
	// if errColnames != nil {
	// 	return nil, errColnames
	// }
	// table.Colnames = &tableColnames

	return &table, nil
}

// ReadTableColnames returns the models.TableColnames
func (o *MyDatastore) ReadTableColnames(t *models.Table, lang string) (*models.TableColnames, error) {

	tableColnames := models.NewColnames(t, lang, nil) // lang=default if empty
	sqlstr, errParam := tableColnames.GetSelectTable()

	fmt.Println("ReadTableColnames, SQL: ", sqlstr)

	if errParam != nil {
		return nil, errParam
	}
	rows, err := o.db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, errCols := rows.Columns()
	if errCols != nil {
		return nil, errCols
	}
	count := len(columns)

	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)

	tableColnames.Header = make([]string, count-1) // the Lang field is apart
	i := 0
	// just 1 time or 0 if lang does not exist
	for rows.Next() {
		for ; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		// discarding the lang field
		for i = 1; i < count; i++ {
			tableColnames.Header[i-1] = string(values[i].([]byte))
		}
	}
	if i == 0 {
		return nil, fmt.Errorf("colnames do not exist for lang param: %s", lang)
	}

	return tableColnames, nil
}

// ReadTableValues returns the models.TableValues
func (o *MyDatastore) ReadTableValues(t *models.Table, start int, count int) (*models.TableValues, error) {
	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}
	tableValues := models.NewValues(t, start, count, rows)

	return tableValues, nil
}
