package dataaccess

import (
	"dataaccess/impl/mydatastore"
	"fmt"
)

func NewDatastore() (IDatastore, error) {
	dbcfg := mydatastore.DBConfig{
		Uid:    "root",
		Pwd:    "secr3tZ",
		IP:     "127.0.0.1",
		Port:   "3306",
		Dbname: "dcgw",
	}

	fmt.Println("Connecting to...", dbcfg.String())

	myDatastore, err := mydatastore.NewDatastore(dbcfg)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return myDatastore, nil
}
