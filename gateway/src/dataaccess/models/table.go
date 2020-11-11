package models

import (
	"fmt"
)

const (
	StatusDeleted = 0 // not more available
	StatusDraft   = 1 // default, available only to owner
	StatusEnabled = 2 // public, available to shoppers
	end           = 3
)

func IsValidStatus(value int) bool {
	return value < int(end)
}

type Table struct {
	Id       string `json:"id"`
	Name     string `json:"name" xml:"name"`
	Owner    string `json:"owner"`
	Descr    string `json:"descr"`
	Tags     string `json:"tags"`
	DefLang  string `json:"deflang"`
	NCols    int    `json:"ncols"`
	NRows    int64  `json:"nrows"`
	Status   int    `json:"status"`
	Colnames *TableColnames
	Values   *TableValues
}

func (o *Table) IsEmpty() bool {
	return o.Owner == "" &&
		o.Name == "" &&
		o.Descr == "" &&
		o.Tags == ""
}

func (o *Table) IsValid() bool {
	return o.Owner != "" &&
		o.Name != "" &&
		o.Descr != "" &&
		o.NCols >= 0 &&
		o.NRows >= 0
}

func (o *Table) String() string {
	return fmt.Sprintf("'%s','%s','%s','%s','%s',%d,%d,%d", o.DefLang, o.Owner, o.Name, o.Descr, o.Tags, o.NCols, o.NRows, o.Status)
}

// func (o *Table) IsOwner() bool {
// 	return o.user.userName == o.Owner
// }

// func (o *Table) SetUser(u *account) bool {
// 	return o.user = u
// }
