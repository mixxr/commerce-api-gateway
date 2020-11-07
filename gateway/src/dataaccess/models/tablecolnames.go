package models

import (
	"fmt"
	"strings"
)

type TableColnames struct {
	table_ *Table
	Lang   string   `json:"lang"`
	Header []string `json:"header"`
}

func (o *TableColnames) String() string {
	return fmt.Sprintf("%s,%s", o.Lang, strings.Join(o.Header, ","))
}

func (o *TableColnames) Parent() *Table {
	return o.table_
}

func (o *TableColnames) SetParent(t *Table) {
	o.table_ = t
}

func (o *TableColnames) IsValid() bool {
	return o.Lang != "" &&
		len(o.Header) > 0 &&
		o.table_ != nil
}

func NewColnames(t *Table, lang string, header []string) *TableColnames {
	newObj := new(TableColnames)
	newObj.table_ = t
	if lang == "" {
		newObj.Lang = t.DefLang
	} else {
		newObj.Lang = lang
	}
	if header != nil {
		newObj.Header = header
		t.NCols = len(header)
	}
	return newObj
}
