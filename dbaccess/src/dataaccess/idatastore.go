package dataaccess

import "dataaccess/models"

type IDatastore interface {
	StoreTable(*models.Table) error
	Read() (*models.Table, error)
}
