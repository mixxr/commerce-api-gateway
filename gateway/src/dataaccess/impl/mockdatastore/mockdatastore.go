package mockdatastore

import (
	"fmt"
	"main/dataaccess/models"
	"math/rand"
	"time"
)

type DBConfig struct {
}

type MockDatastore struct {
	rndGen *rand.Rand
	m      map[string]*models.Table
}

func (o DBConfig) String() string {
	return fmt.Sprintf("mockdatastore")
}

func NewDatastore(o *DBConfig) (*MockDatastore, error) {
	newObj := new(MockDatastore)
	s1 := rand.NewSource(time.Now().UnixNano())
	newObj.rndGen = rand.New(s1)
	newObj.m = make(map[string]*models.Table)

	return newObj, nil
}

// BEGIN Store functions

func (o *MockDatastore) StoreTable(tin *models.Table) error {
	if o.m[tin.Owner] != nil {
		return fmt.Errorf("Error 1062: Duplicate entry '%s-%s' for key 'owner'", tin.Owner, tin.Name)
	}
	o.m[tin.Owner] = tin
	return nil
}

func (o *MockDatastore) StoreTableColnames(t *models.TableColnames) error {
	return nil
}

func (o *MockDatastore) StoreTableValues(t *models.TableValues) error {
	return nil
}

func (o *MockDatastore) UpdateTable(t *models.Table) error {
	if o.m[t.Owner] == nil {
		return fmt.Errorf("service do not exist for params: %s", t.String())
	}
	o.m[t.Owner] = t
	return nil
}

// func (o *MockDatastore) AddColnames(t *models.TableColnames) error {
// 	return nil
// }

// func (o *MockDatastore) AddValues(t *models.TableValues) error {

// 	return nil
// }

// END Strore functions
func (o *MockDatastore) ReadTables(tin *models.Table) ([]*models.Table, error) {
	return nil, fmt.Errorf("TBD")
}

func (o *MockDatastore) ReadTable(tin *models.Table) (*models.Table, error) {

	t := o.m[tin.Owner]

	if t != nil && t.Status >= tin.Status {
		return t, nil
	}

	return nil, fmt.Errorf("service do not exist for params: %s", tin.String())
}

// ReadTableColnames returns the models.TableColnames
func (o *MockDatastore) ReadTableColnames(t *models.Table, lang string) (*models.TableColnames, error) {
	headers := []string{"marca", "modello", "prezzo", "valuta"}
	tableColnames := models.NewColnames(t, lang, headers)

	return tableColnames, nil
}

// ReadTableValues returns the models.TableValues
func (o *MockDatastore) ReadTableValues(t *models.Table, start int, count int64) (*models.TableValues, error) {
	rows := [][]string{
		{"fiat", "uno 1.0 fire", "5.000", "lire"},
		{"fiat", "uno 1.4 TD", "10.000", "lire"},
		{"fiat", "panda 750 fire", "4.000", "lire"},
		{"fiat", "127 900", "4.500", "lire"},
		{"fiat", "128 1.2", "5.500", "lire"},
	}
	tableValues := models.NewValues(t, start, count, rows)

	return tableValues, nil
}

func (o *MockDatastore) DeleteTable(*models.Table) error {
	return nil
}

func (o *MockDatastore) DeleteTableColnames(t *models.Table, langs []string) error {
	return nil
}

func (o *MockDatastore) DeleteTableValues(t *models.Table, count int64) error {
	return nil
}
