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

// StoreTable in a Transaction:
// 1. insert into table_
// 2. create table owner_name_cols
// 3. create table owner_name_values
// 4. insert into owner_name_cols
// 5. insert int owner_name_values
func (o *MyDatastore) StoreTable(t *models.Table) error {
	var sqlstr [5]string
	var err error
	sqlstr[0], err = t.GetInsertTable()
	if err != nil {
		return err
	}
	sqlstr[1], err = t.Colnames.GetCreateTable()
	if err != nil {
		return err
	}
	sqlstr[2], err = t.Values.GetCreateTable()
	if err != nil {
		return err
	}
	sqlstr[3], err = t.Colnames.GetInsertTable()
	if err != nil {
		return err
	}
	sqlstr[4], err = t.Values.GetInsertTable()
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

// 1. update table_ (descr, tags)
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

func (o *MyDatastore) Read() (*models.Table, error) {

	table := models.Table{
		Name:    "example",
		Owner:   "michelangelo190283",
		DefLang: "it",
		Tags:    "tag1,tag2",
		Descr:   "auto FIAT anni 80",
	}

	headers := []string{"marca", "modello", "prezzo", "valuta"}

	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}

	// var cols *models.TableColnames = models.NewColnames(&table, &headers)
	// var values *models.TableValues = models.NewValues(&table, &rows)

	table.Colnames = models.NewColnames(&table, headers)
	table.Values = models.NewValues(&table, rows)

	return &table, nil
}
