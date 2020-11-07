package mydatastore

import (
	"dataaccess/impl/mydatastore/utils"
	"dataaccess/models"
	"database/sql"
	"fmt"
	"log"
	"os"
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
	if mainconn == nil {
		// log
		file, err := os.OpenFile("dataaccess.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		log.SetOutput(file)

		log.Println("MyDataStore - SQL DB!")
	}
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
	var version string
	var tot int64
	db.QueryRow("SELECT VERSION()").Scan(&version)
	log.Println("Connected to:", version)

	db.QueryRow("SELECT count(*) as tot from table_").Scan(&tot)
	log.Println("DB contains tables:", tot)

	mainconn = db

	return db, nil
}

// BEGIN Store functions

// StoreTable adds 1 row as
// 1. insert into table_
func (o *MyDatastore) StoreTable(t *models.Table) error {
	if !t.IsValid() {
		return fmt.Errorf("service cannot be created for this params: %s", t.String())
	}
	var sqlstr string
	var err error
	sqlstr, err = utils.GetInsertTable(t)
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
	if !t.IsValid() {
		return fmt.Errorf("colnames cannot be created for this params: %s", t.String())
	}
	var sqlstr [3]string
	var err error

	t.Parent().NCols = len(t.Header)

	sqlstr[0], err = utils.GetCreateTableColnames(t)
	if err != nil {
		return err
	}
	sqlstr[1], err = utils.GetInsertTableColnames(t)
	if err != nil {
		return err
	}
	sqlstr[2], err = utils.GetUpdateNCols(t.Parent())
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
// 3. update table_ with GetIncrementTable(affectedRows)
func (o *MyDatastore) StoreTableValues(t *models.TableValues) error {
	if !t.IsValid() {
		return fmt.Errorf("values cannot be stored for this params: %s", t.Parent())
	}
	var sqlstr [2]string
	var err error

	sqlstr[0], err = utils.GetCreateTableValues(t)
	if err != nil {
		return err
	}
	sqlstr[1], err = utils.GetInsertTableValues(t)
	if err != nil {
		return err
	}

	// Transaction starts
	tx, err := o.db.Begin()
	if err != nil {
		return err
	}
	t.Count = 0

	// CREATE TABLE IF DOES NOT EXISTS
	_, err = tx.Exec(sqlstr[0])
	if err != nil {
		tx.Rollback()
		return err
	}
	// INSERT ROWS
	res, errINS := tx.Exec(sqlstr[1])
	if errINS != nil {
		tx.Rollback()
		return errINS
	}
	var affectedRows int64
	affectedRows, errINS = res.RowsAffected()
	if errINS != nil {
		tx.Rollback()
		return errINS
	}
	// UPDATE table_ by affectedRows
	sqlUpdate, errUPD := utils.GetIncrementTable(t.Parent(), affectedRows)
	if errUPD != nil {
		tx.Rollback()
		return errUPD
	}
	_, errUPD = tx.Exec(sqlUpdate)
	if errUPD != nil {
		tx.Rollback()
		return errUPD
	}

	if err = tx.Commit(); err == nil {
		t.Count = affectedRows
	}

	return err
}

// UpdateTable just updates 2 fields of table_: descr, tags
func (o *MyDatastore) UpdateTable(t *models.Table) error {
	sql, err := utils.GetUpdateTable(t)
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
func (o *MyDatastore) AddColnames(t *models.TableColnames) error {
	var sql1, sql2 string
	var err error
	sql1, err = utils.GetDeleteTableColnames(t)
	if err != nil {
		return err
	}
	sql2, err = utils.GetInsertTableColnames(t)
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
func (o *MyDatastore) AddValues(t *models.TableValues) error {
	var sql1, sql2 string
	var err error
	sql1, err = utils.GetInsertTableValues(t)
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

	sql2, err = utils.GetIncrementTable(t.Parent(), nrows)
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

// ReadTables returns  []models.Table without colnames neither values
func (o *MyDatastore) ReadTables(tin *models.Table) ([]*models.Table, error) {
	sqlstr, errParam := utils.GetSelectSearchTable(tin)
	log.Println("ReadTables, input: ", tin, sqlstr)

	if errParam != nil {
		return nil, errParam
	}
	rows, err := o.db.Query(sqlstr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []*models.Table

	for rows.Next() {
		table := models.Table{}
		if err = rows.Scan(&table.Id, &table.Owner, &table.Name, &table.Descr, &table.Tags, &table.DefLang, &table.NCols, &table.NRows); err != nil {
			return nil, err
		}
		log.Println("ReadTables, table: ", table)
		tables = append(tables, &table)
	}

	return tables, nil
}

// ReadTable returns the models.Table without colnames neither values
func (o *MyDatastore) ReadTable(tin *models.Table) (*models.Table, error) {
	name := tin.Name
	owner := tin.Owner
	sqlstr, errParam := utils.GetSelectTable(name, owner)

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
func (o *MyDatastore) ReadTableColnames(tin *models.Table, lang string) (*models.TableColnames, error) {

	log.Println("ReadTableColnames, input: ", tin, lang)

	// TODO: it is needed for NCols that is not part of the input. A different SELECT can be done without using fixed fields
	t, errMain := o.ReadTable(tin)
	if errMain != nil {
		return nil, fmt.Errorf("colnames do not exist for input param: %s", tin.String())
	}

	tableColnames := models.NewColnames(t, lang, nil) // lang=default if empty
	sqlstr, errParam := utils.GetSelectTableColnames(tableColnames)

	log.Println("ReadTableColnames, SQL: ", sqlstr, errParam)

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
		log.Println("colnames do not exist for lang param:", lang)
		return nil, fmt.Errorf("colnames do not exist for lang param: %s", lang)
	}

	return tableColnames, nil
}

// ReadTableValues returns the models.TableValues
func (o *MyDatastore) ReadTableValues(tin *models.Table, start int, count int64) (*models.TableValues, error) {

	log.Println("ReadTableValues, input: ", tin, start, count)

	// TODO: it is needed for NCols that is not part of the input. A different SELECT can be done without using fixed fields
	t, errMain := o.ReadTable(tin)
	if errMain != nil {
		return nil, fmt.Errorf("values do not exist for input param: %s", tin.String())
	}
	tableValues := models.NewValues(t, start, count, nil)
	sqlstr, errParam := utils.GetSelectTableValues(tableValues)

	log.Println("ReadTableValues, SQL: ", sqlstr, errParam)

	if errParam != nil {
		return nil, errParam
	}
	// real tot amount of selected rows
	var totRows int64
	sqlCount, _ := utils.GetSelectNRows(tin)
	o.db.QueryRow(sqlCount).Scan(&totRows)
	fmt.Println("ReadTableValues, SQL, result, count: ", sqlCount, totRows, count)
	if totRows > 0 {
		if totRows > count {
			totRows = count // this allocates always the min(available #rows, requested #rows)
		}
		// values
		rows, err := o.db.Query(sqlstr)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		columns, errCols := rows.Columns()
		if errCols != nil {
			return nil, errCols
		}
		totCols := len(columns)

		values := make([]interface{}, totCols)
		valuePtrs := make([]interface{}, totCols)

		tableValues.Rows = make([][]string, totRows)
		for r := int64(0); r < totRows; r++ {
			tableValues.Rows[r] = make([]string, totCols)
		}
		j := 0
		for rows.Next() {
			for i := 0; i < totCols; i++ {
				valuePtrs[i] = &values[i]
			}
			rows.Scan(valuePtrs...)
			for i := 0; i < totCols; i++ {
				tableValues.Rows[j][i] = string(values[i].([]byte))
			}
			j++
		}
		tableValues.Count = totRows

		return tableValues, nil

	}
	log.Println("values do not exist for params:", tin)
	return nil, fmt.Errorf("values do not exist for params: %s", tin)

}

// DELETE functions

// DeleteTable in a Transaction:
// 1. delete the table_ row
// 2. DROP owner_name_colnames
// 3. DROP owner_name_values
func (o *MyDatastore) DeleteTable(t *models.Table) error {
	var sqlstr [2]string
	var err error

	sqlstr[0], err = utils.GetDeleteTable(t)
	if err != nil {
		return err
	}
	sqlstr[1], err = utils.GetDropTables(t)
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
