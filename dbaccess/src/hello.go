package main

import (
	"dataaccess"
	"dataaccess/impl/mockdatastore"
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

func prepareMock() *mockdatastore.MockDatastore {

	dbcfg := mockdatastore.DBConfig{}

	fmt.Println("Connecting to..." + dbcfg.String())

	mockDatastore, err := mockdatastore.NewDatastore(dbcfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return mockDatastore
}

func main() {
	fmt.Printf("Starting Main...\n")
	var err error
	var table1 *models.Table
	var table2 *models.Table

	// Mock DS
	var mockDS dataaccess.IDatastore
	mockDS = prepareMock()
	table1, err = mockDS.ReadTable("service", "micser")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Read Mock:", table1.String())

	// MySQL DS
	var myDS dataaccess.IDatastore
	myDS = prepareMySQL()

	err = myDS.StoreTable(table1)
	if err != nil {
		fmt.Println("Store SQL Error:", err)
	}
	fmt.Println("Store SQL:", table1)

	name, owner := table1.Name, table1.Owner
	table2, err = myDS.ReadTable(name, owner)

	if err != nil {
		fmt.Println("Read SQL Error:", err)
	} else {
		fmt.Println("Read SQL:", table2.String())
	}

	fmt.Println("====END")
}
