package models

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"io"
	"strings"
)

type TableValues struct {
	table_ *Table
	Start  int
	Count  int64
	Rows   [][]string
}

func NewValues(t *Table, start int, count int64, rows [][]string) *TableValues {
	newObj := new(TableValues)
	newObj.table_ = t
	newObj.Start = start
	newObj.Count = count
	if rows != nil {
		newObj.Rows = rows
		//newObj.Count = len(rows) cannot assign, Count is used as control field upon SQL operation
	}
	return newObj
}

func (o *TableValues) Size() int {
	return len(o.Rows) * len(o.Rows[0])
}

func (o *TableValues) String() string {
	var buffer bytes.Buffer

	for _, row := range o.Rows {
		// TODO: escape comma for CSV
		str := strings.Join(row, ",")
		fmt.Fprintf(&buffer, "%s\n", str)
	}

	return buffer.String()
}

func (o *TableValues) StreamCSV(writer io.Writer) {
	w := csv.NewWriter(writer)

	// for _, row := range o.Rows {
	// 	w.Write(row)
	// 	//		fmt.Fprintf(&buffer, "%s\n", str)
	// }
	w.WriteAll(o.Rows)
}

func (o *TableValues) SetParent(t *Table) {
	o.table_ = t
}

func (o *TableValues) Parent() *Table {
	return o.table_
}

func (o *TableValues) IsValid() bool {
	return o.Start >= 0 &&
		len(o.Rows) > 0 &&
		o.table_ != nil
}
