package dataaccess

import "dataaccess/models"

type IDatastore interface {
	StoreTable(*models.Table) error
	ReadTable(name string, owner string) (*models.Table, error)
}
