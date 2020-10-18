package models

import (
	"fmt"
	"strings"
)

type Row struct {
	Col1 string
	Col2 string
	Col3 string
	Col4 string
	Col5 string
}

type Table struct {
	Id    string
	Name  string
	Owner string
	Tags  []string
	Start int
	Rows  []Row
}

func (o *Table) String() string {
	return fmt.Sprintf("%s;%s;%s;%s;%d;%d", o.Id, o.Name, o.Owner, strings.Join(o.Tags, ","), o.Start, len(o.Rows))
}
