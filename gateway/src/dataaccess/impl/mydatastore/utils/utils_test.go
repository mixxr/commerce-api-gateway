package utils_test

import (
	"fmt"
	"main/dataaccess/impl/mydatastore/utils"
	"main/dataaccess/models"
	"testing"
)

var table *models.Table
var tableColnames *models.TableColnames
var tableValues *models.TableValues

var lang, name, owner string
var start, count int

func prepareData() {
	lang = "en"
	start = 0
	count = 5
	name = "name"
	owner = "owner"

	table = &models.Table{
		Name:    name,
		Owner:   owner,
		DefLang: "it",
		Tags:    "tag1,tag2",
		Descr:   "short description about the service",
	}
	headers := []string{"marca", "modello", "prezzo", "valuta"}
	tableColnames = models.NewColnames(table, lang, headers)

	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}
	tableValues = models.NewValues(table, start, int64(count), rows)

}

func TestUtilsCreateTables(t *testing.T) {
	prepareData()

	var sqlstr [3]string
	var sqlstrOK [3]string

	sqlstrOK[0] = fmt.Sprintf(`UPDATE table_ SET nrows=nrows+%d WHERE owner='%s' AND name='%s';`,
		count,
		owner,
		name)
	sqlstrOK[1] = fmt.Sprintf(`INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES(%s);`,
		table.String())
	sqlstrOK[2] = fmt.Sprintf(`UPDATE table_ SET descr='%s',tags='%s',ncols=%d,nrows=%d WHERE owner='%s' AND name='%s';`,
		table.Descr,
		table.Tags,
		table.NCols,
		table.NRows,
		table.Owner,
		table.Name)

	sqlstr[0], _ = utils.GetIncrementTable(table, int64(count))
	sqlstr[1], _ = utils.GetInsertTable(table)
	sqlstr[2], _ = utils.GetUpdateTable(table)
	// sqlstr[3] = utils.GetSelectTable(name, owner)
	// sqlstr[4] = utils.GetUpdateNCols(table)

	for i, stmt := range sqlstr {
		if stmt != sqlstr[i] {
			t.Errorf("SQL not correct, got %s, want %s", stmt, sqlstr[i])
		}
	}
}
