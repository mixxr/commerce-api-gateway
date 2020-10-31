package models

import (
	"fmt"
)

type Table struct {
	Id       string `json:"id"`
	Name     string `json:"name" xml:"name"`
	Owner    string `json:"owner"`
	Descr    string `json:"descr"`
	Tags     string `json:"tags"`
	DefLang  string `json:"deflang"`
	NCols    int    `json:"ncols"`
	NRows    int    `json:"nrows"`
	Colnames *TableColnames
	Values   *TableValues
}

func (o *Table) String() string {
	return fmt.Sprintf("'%s','%s','%s','%s','%s',%d,%d", o.DefLang, o.Owner, o.Name, o.Descr, o.Tags, o.NCols, o.NRows)
}
