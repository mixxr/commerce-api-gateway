package impl

import (
	"dataaccess/models"
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type DBConfig struct {
	Uid    string
	Pwd    string
	IP     string
	Port   string
	Dbname string
}

type DB struct {
	*sql.DB
}

func (o DBConfig) String() string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", o.Uid, o.Pwd, o.IP, o.Port, o.Dbname)
}

func (o DBConfig) Connect() (*DB, error) {
	db, err := sql.Open("mysql", o.String())
	if err != nil {
		return nil, err
	}
	//defer db.Close()

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Connect and check the server version
	var version string
	db.QueryRow("SELECT VERSION()").Scan(&version)
	fmt.Println("Connected to:", version)

	return &DB{db}, nil
}

func (db *DB) Store(*models.Table) (bool, error) {
	return true, nil
}

func (db *DB) Read() (*models.Table, error) {
	rows := []models.Row{
		{Col1: "ciao"},
		{Col1: "ciao", Col2: "buongiorno"},
	}
	table := models.Table{
		Id:    "349DU30U84DD34URF8",
		Name:  "Example",
		Start: 0,
		Tags:  []string{"tag1,tag2"},
		Owner: "michelangelo190283",
		Rows:  rows,
	}
	return &table, nil
}
