package models

import (
	"fmt"
)

type TableValues struct {
	table_ *Table
	Start  int
	Count  int
	Rows   [][]string
}

func NewValues(t *Table, start int, count int, rows [][]string) *TableValues {
	newObj := new(TableValues)
	newObj.table_ = t
	newObj.Start = start
	newObj.Count = count
	if rows != nil {
		newObj.Rows = rows
		newObj.Count = len(rows)
	}
	return newObj
}

func (o *TableValues) Size() int {
	return len(o.Rows) * len(o.Rows[0])
}

func (o *TableValues) String() string {
	return fmt.Sprintf("%d,%d", o.Start, o.Count)
}

func (o *TableValues) Parent() *Table {
	return o.table_
}
