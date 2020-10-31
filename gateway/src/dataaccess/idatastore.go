package dataaccess

import "dataaccess/models"

type IDatastore interface {
	// create
	StoreTable(*models.Table) error
	StoreTableColnames(t *models.TableColnames) error
	StoreTableValues(t *models.TableValues) error
	// update
	UpdateTable(t *models.Table) error         // updates table_
	AddColnames(t *models.TableColnames) error // add additional entry in <table>_colnames
	AddValues(t *models.TableValues) error     // add additional entries in <table>_values
	// read
	ReadTable(name string, owner string) (*models.Table, error)
	ReadTableColnames(o *models.Table, lang string) (*models.TableColnames, error)
	ReadTableValues(t *models.Table, start int, count int) (*models.TableValues, error)
	// delete
	// TODO: delete functions
}
