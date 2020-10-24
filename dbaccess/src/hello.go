package main

import (
	"dataaccess"
	"dataaccess/impl"
	"dataaccess/models"
	"fmt"
	"os"
)

type Env struct {
	db dataaccess.IDatastore
}

func prepareMySQL() *Env {
	fmt.Printf("Starting Main...\n")
	dbcfg := impl.DBConfig{
		Uid:    "root",
		Pwd:    "secr3tZ",
		IP:     "127.0.0.1",
		Port:   "3306",
		Dbname: "dcgw",
	}

	dbcfg.String()
	var db *impl.DB
	var err error
	db, err = dbcfg.Connect()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return &Env{db}
}

func main() {

	env := prepareMySQL()
	//table := models.Table
	var err error
	var result bool
	result, err = env.db.Store(nil)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Stored:", result)

	var table *models.Table
	table, err = env.db.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Read:", table.String())

	s, _ := table.Cols.GetCreateTable()
	fmt.Println("GetCreateTable:", s)
	s, _ = table.Cols.GetInsertTable()
	fmt.Println("GetInsertTable:", s)
	s, _ = table.Values.GetCreateTable()
	fmt.Println("GetCreateTable:", s)
	s, _ = table.Values.GetInsertTable()
	fmt.Println("GetInsertTable:", s)
}
