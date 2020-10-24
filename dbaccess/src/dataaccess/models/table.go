package models

import (
	"fmt"
)

type Table struct {
	Id      string `json:"id"`
	Name    string `json:"name"`
	Owner   string `json:"owner"`
	Descr   string `json:"descr"`
	Tags    string `json:"tags"`
	DefLang string `json:"deflang"`
	NCols   int    `json:"ncols"`
	NRows   int    `json:"nrows"`
	Start   int
	Cols    *TableCols
	Values  *TableValues
}

func (o *Table) String() string {
	return fmt.Sprintf("%s;%s;%s;%s;%d;%d", o.Id, o.Name, o.Owner, o.Tags, o.Start, o.Values.Size())
}
