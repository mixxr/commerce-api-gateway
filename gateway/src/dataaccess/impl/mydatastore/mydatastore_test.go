package mydatastore_test

import (
	"dataaccess"
	"dataaccess/impl/mydatastore"
	"dataaccess/models"
	"fmt"
	"logger"
	"math/rand"
	"os"
	"testing"
	"time"
)

var myDatastore *mydatastore.MyDatastore
var file *os.File
var name, owner string
var ncols int

func prepareMySQL() *mydatastore.MyDatastore {

	if myDatastore != nil {
		return myDatastore
	}

	logger.AppLogger.Info("test", "test", "MyDataStoreTest - initalizing...")

	rand.Seed(time.Now().UnixNano())
	name = fmt.Sprintf("service_%d", rand.Intn(1000))
	owner = fmt.Sprintf("owner_%d", rand.Intn(1000))
	ncols = 4
	logger.AppLogger.Info("test", "test", "Test Parameters...", name, owner, ncols)

	dbcfg := mydatastore.DBConfig{
		Uid:    "root",
		Pwd:    "secr3tZ",
		IP:     "127.0.0.1",
		Port:   "3306",
		Dbname: "dcgw",
	}

	logger.AppLogger.Info("test", "test", "Connecting to...", dbcfg.String())
	var err error
	myDatastore, err = mydatastore.NewDatastore(&dbcfg)
	if err != nil {
		logger.AppLogger.Fatal("test", "test", err)
	}

	return myDatastore
}

func TestMyDatastoreCreate(t *testing.T) {
	var err error

	// MySQL DS
	var myDS dataaccess.IDatastore
	myDS = prepareMySQL()

	// table_
	table1 := models.Table{Name: name, Owner: owner, DefLang: "it", Descr: "test service", Tags: "test,dummy,golang"}
	err = myDS.StoreTable(&table1)
	if err != nil {
		t.Errorf("StoreTable ERROR: %s", err)
	}
	logger.AppLogger.Info("test", "test", "StoreTable SQL:", table1)

	// tablecolnames
	header := []string{"test_colname1", "test_colname2", "test_colname3", "test_colname4"}
	tableCNs := models.NewColnames(&table1, "", header)

	err = myDS.StoreTableColnames(tableCNs)
	if err != nil {
		t.Errorf("StoreTableColnames ERROR: %s", err)
	}
	logger.AppLogger.Info("test", "test", "StoreTableColnames SQL:", tableCNs)

	// tablevalues
	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}
	start := 100
	count := 5
	tableVs := models.NewValues(&table1, start, int64(count), rows)
	err = myDS.StoreTableValues(tableVs)
	if err != nil {
		t.Errorf("StoreTableValues ERROR: %s", err)
	}
	logger.AppLogger.Info("test", "test", "StoreTableValues SQL:", tableVs)
}

func TestMyDatastoreRead(t *testing.T) {
	var err error
	var table1 *models.Table
	//	rndGen := new(rand.Rand)

	// MySQL DS
	var myDS dataaccess.IDatastore
	myDS = prepareMySQL()

	// table_
	tin := models.Table{Name: name, Owner: owner}
	table1, err = myDS.ReadTable(&tin)
	if err != nil {
		t.Errorf("ReadTable Error: %s", err)
	}
	if table1.Name != name {
		t.Errorf("table.Name is not correct, got %s, want %s", table1.Name, name)
	}
	if table1.Owner != owner {
		t.Errorf("table.Owner is not correct, got %s, want %s", table1.Owner, owner)
	}
	if table1.NCols != ncols {
		t.Errorf("table.NCols is not correct, got %d, want %d", table1.NCols, ncols)
	}
	logger.AppLogger.Info("test", "test", "table_ READ:", table1)

	// colnames
	var tableColnames *models.TableColnames
	tableColnames, err = myDS.ReadTableColnames(table1, "")
	if err != nil {
		t.Errorf("ReadTableColnames Error: %s", err)
	}
	logger.AppLogger.Info("test", "test", "ReadTableColnames SQL:", tableColnames)
	testColname := fmt.Sprintf("test_colname%d", table1.NCols)
	if tableColnames.Lang != table1.DefLang {
		t.Errorf("ReadTableColnames deflang is not correct, got %s, want %s", tableColnames.Lang, table1.DefLang)
	}
	if len(tableColnames.Header) != table1.NCols {
		t.Errorf("ReadTableColnames NCols is not correct, got %d, want %d", len(tableColnames.Header), table1.NCols)
	}
	if tableColnames.Header[table1.NCols-1] != testColname {
		t.Errorf("ReadTableColnames testColname is not correct, got %s, want %s", tableColnames.Header[table1.NCols-1], testColname)
	}

	// values
	var tableValues *models.TableValues
	start, count := 0, 50
	tableValues, err = myDS.ReadTableValues(table1, start, int64(count))
	if err != nil {
		t.Errorf("ReadTableValues Error: %s", err)
	} else {
		logger.AppLogger.Info("test", "test", "ReadTableValues SQL:", tableValues)

		if tableValues.Start != start {
			t.Errorf("ReadTableValues Start is not correct, got %d, want %d", tableValues.Start, start)
		}
		if len(tableValues.Rows) > count || int64(len(tableValues.Rows)) != tableValues.Count {
			t.Errorf("ReadTableValues count with %d as param is not correct, got %d, want %d", count, len(tableValues.Rows), tableValues.Count)
		}
		if tableValues.Rows[tableValues.Count-1][table1.NCols-1] != "lire" {
			t.Errorf("ReadTableValues Rows[last][last] is not correct, got %s, want %s", tableValues.Rows[tableValues.Count-1][table1.NCols-1], "lire")
		}
	}

}

func TestMyDatastoreRemove(t *testing.T) {
	var err error
	var table1 *models.Table
	//	rndGen := new(rand.Rand)

	// MySQL DS
	var myDS dataaccess.IDatastore
	myDS = prepareMySQL()

	// table_
	tin := models.Table{Name: name, Owner: owner}
	err = myDS.DeleteTable(&tin)
	if err != nil {
		t.Errorf("DeleteTable Error: %s", err)
	} else {
		logger.AppLogger.Info("test", "test", "DeleteTable OK:", table1)
	}
}
