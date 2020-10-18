package dataaccess

import "dataaccess/models"

type IDatastore interface {
	Store(*models.Table) (bool, error)
	Read() (*models.Table, error)
}
