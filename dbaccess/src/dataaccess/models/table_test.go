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

func TestInitializations(t *testing.T) {
	// var err error
	//	rndGen := new(rand.Rand)

	table1 := models.Table{Name: "testservice", Owner: "testowner", DefLang: "it", Descr: "test service", Tags: "test,dummy,golang"}

	header := []string{"test_colname1", "test_colname2", "test_colname3", "test_colname4"} // colnames
	tableCNs := models.NewColnames(&table1, "", header)

	if tableCNs.Lang != table1.DefLang {
		t.Errorf("table.deflang is not propagated correctly, got %s, want %s", tableCNs.Lang, table1.DefLang)
	}
	if tableCNs.Parent() != &table1 {
		t.Errorf("table is not propagated correctly, got %s, want %s", tableCNs.Parent(), &table1)
	}
	if len(header) != table1.NCols {
		t.Errorf("table.NCols is not calculated correctly, got %d, want %d", table1.NCols, len(header))
	}

	tableCNs2 := models.NewColnames(&table1, "de", header)
	if tableCNs2.Lang != "de" {
		t.Errorf("table.deflang is not initialized correctly, got %s, want %s", tableCNs.Lang, "de")
	}

	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}
	start := 100
	count := 5
	size := 20
	tableVs := models.NewValues(&table1, start, count, rows)

	if tableVs.Count != len(rows) {
		t.Errorf("tablevalues.Count is not calculated correctly, got %d, want %d", tableVs.Count, len(rows))
	}
	if tableVs.Parent() != &table1 {
		t.Errorf("table is not propagated correctly, got %s, want %s", tableVs.Parent(), &table1)
	}
	if tableVs.Start != start {
		t.Errorf("table.NCols is not calculated correctly, got %d, want %d", tableVs.Start, start)
	}
	if tableVs.Size() != size {
		t.Errorf("table.Size() is not calculated correctly, got %d, want %d", tableVs.Size(), size)
	}
}
