package main

import (
	"dataaccess"
	"dataaccess/impl/mydatastore"
	"dataaccess/models"
	"fmt"
	"os"
)

// type Env struct {
// 	db dataaccess.IDatastore
// }

func prepareMySQL() *mydatastore.MyDatastore {

	dbcfg := mydatastore.DBConfig{
		Uid:    "root",
		Pwd:    "secr3tZ",
		IP:     "127.0.0.1",
		Port:   "3306",
		Dbname: "dcgw",
	}

	fmt.Println("Connecting to..." + dbcfg.String())

	myDatastore, err := mydatastore.NewDatastore(dbcfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return myDatastore
}

func main() {
	fmt.Printf("Starting Main...\n")

	var myDS dataaccess.IDatastore
	myDS = prepareMySQL()
	//table := models.Table
	var err error

	var table *models.Table
	table, err = myDS.Read()
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Read:", table.String())

	s, _ := table.Colnames.GetCreateTable()
	fmt.Println("GetCreateTable:", s)
	s, _ = table.Colnames.GetInsertTable()
	fmt.Println("GetInsertTable:", s)
	s, _ = table.Values.GetCreateTable()
	fmt.Println("GetCreateTable:", s)
	s, _ = table.Values.GetInsertTable()
	fmt.Println("GetInsertTable:", s)

	err = myDS.StoreTable(table)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("====END")
}
