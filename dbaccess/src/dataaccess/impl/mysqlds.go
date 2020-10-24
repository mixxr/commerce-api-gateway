package impl

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

type DB struct {
	*sql.DB
}

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

func (o DBConfig) Connect() (*DB, error) {
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

	return &DB{db}, nil
}

func (db *DB) Store(*models.Table) (bool, error) {
	return true, nil
}

func (db *DB) Read() (*models.Table, error) {

	table := models.Table{
		Id:      "349DU30U84DD34URF8",
		Name:    "example",
		Start:   0,
		Tags:    "tag1,tag2",
		Owner:   "michelangelo190283",
		DefLang: "it",
	}

	headers := []string{"marca", "modello", "prezzo", "valuta"}

	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}

	// var cols *models.TableCols = models.NewCols(&table, &headers)
	// var values *models.TableValues = models.NewValues(&table, &rows)

	table.Cols = models.NewCols(&table, headers)
	table.Values = models.NewValues(&table, rows)

	return &table, nil
}
