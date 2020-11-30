package dataaccess

import "main/dataaccess/models"

type IDatastore interface {
	// create
	StoreTable(*models.Table) error
	StoreTableColnames(t *models.TableColnames) error
	StoreTableValues(t *models.TableValues) error
	// update
	UpdateTable(t *models.Table) error // updates table_
	//AddColnames(t *models.TableColnames) error // add additional entry in <table>_colnames
	//AddValues(t *models.TableValues) error // add additional entries in <table>_values
	// read
	ReadTables(*models.Table) ([]*models.Table, error)
	ReadTable(*models.Table) (*models.Table, error)
	ReadTableColnames(o *models.Table, lang string) (*models.TableColnames, error)
	ReadTableValues(t *models.Table, start int, count int64) (*models.TableValues, error)
	// delete
	// TODO: delete functions
	DeleteTable(*models.Table) error
	DeleteTableColnames(*models.Table, []string) error
	DeleteTableValues(*models.Table, int64) error
}
