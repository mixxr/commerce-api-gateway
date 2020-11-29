package mockdatastore_test

import (
	"fmt"
	"main/dataaccess"
	"main/dataaccess/impl/mockdatastore"
	"main/dataaccess/models"
	"os"
	"testing"
)

func prepareMock() *mockdatastore.MockDatastore {

	dbcfg := &mockdatastore.DBConfig{}

	fmt.Println("Connecting to..." + dbcfg.String())

	mockDatastore, err := mockdatastore.NewDatastore(dbcfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	return mockDatastore
}

func TestMockDatastore(t *testing.T) {
	var err error
	var table1 *models.Table
	var mockDS dataaccess.IDatastore
	mockDS = prepareMock()
	tin := models.Table{Name: "service", Owner: "micser"}
	table1, err = mockDS.ReadTable(&tin)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Read Mock:", table1.String())

		table1.Colnames, _ = mockDS.ReadTableColnames(table1, "")
		table1.Values, _ = mockDS.ReadTableValues(table1, 0, 5)
	}

}
