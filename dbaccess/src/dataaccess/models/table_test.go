package models_test

import (
	"dataaccess/models"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"testing"
)

func TestTable(t *testing.T) {
	jsonFilePath := "examples/table.json"
	jsonFile, err := os.Open(jsonFilePath)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + jsonFilePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	var table models.Table
	json.Unmarshal(byteValue, &table)

	//  {
	//     "id": "",
	//     "name":"bicycle",
	//     "owner": "sam",
	//     "deflang": "en",
	//     "descr": "bicycle models for summer 2020",
	//     "tags": "summer,2020"
	// }
	if table.Name != "bicycle" {
		t.Errorf("table.name, got %s, want %s", table.Name, "bicycle")
	}
	if table.Owner != "sam" {
		t.Errorf("table.owner, got %s, want %s", table.Owner, "sam")
	}

}

func TestTableColnames(t *testing.T) {
	jsonFilePath := "examples/tablecols.json"
	jsonFile, err := os.Open(jsonFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + jsonFilePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var tablecols models.TableColnames

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'table' which we defined above
	json.Unmarshal(byteValue, &tablecols)

	//  {
	// "lang": "en",
	// "header": ["code","color","size","price","currency"],
	// }
	if tablecols.Lang != "en" {
		t.Errorf("tablecols.Lang, got %s, want %s", tablecols.Lang, "en")
	}
	if tablecols.Header[0] != "code" {
		t.Errorf("tablecols.Header[0], got %s, want %s", tablecols.Header[0], "code")
	}

}

func TestTableValues(t *testing.T) {
	jsonFilePath := "examples/tablevalues.json"
	jsonFile, err := os.Open(jsonFilePath)
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully Opened " + jsonFilePath)
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Users array
	var tablevalues models.TableValues

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'table' which we defined above
	json.Unmarshal(byteValue, &tablevalues)

	//  {
	//     "rows": [
	//         ["AKL2374","RED","XS","123","USD"],
	//         ["QWE56712","WHITE","XS","120","USD"],
	//         ["TRF82374","BLACK","XS","120","USD"],
	//         ["RTG23734","RED","M","125","USD"],
	//         ["QAD2375","RED","XS","125","USD"],
	//         ["TYH2346","RED","L","125","USD"]
	//    ]
	// }
	if len(tablevalues.Rows) != 6 {
		t.Errorf("tablevalues.Rows Len, got %d, want %d", len(tablevalues.Rows), 6)
	}
	if tablevalues.Rows[0][4] != "USD" {
		t.Errorf("tablevalues.Rows[0][4], got %s, want %s", tablevalues.Rows[0][4], "USD")
	}

}
