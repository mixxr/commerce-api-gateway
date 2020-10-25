package models

import (
	"fmt"
)

type Table struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Owner    string `json:"owner"`
	Descr    string `json:"descr"`
	Tags     string `json:"tags"`
	DefLang  string `json:"deflang"`
	NCols    int    `json:"ncols"`
	NRows    int    `json:"nrows"`
	Start    int
	Colnames *TableColnames
	Values   *TableValues
}

func (o *Table) String() string {
	return fmt.Sprintf("'%s','%s','%s','%s','%s',%d,%d", o.DefLang, o.Owner, o.Name, o.Descr, o.Tags, o.NCols, o.NRows)
}

// GetInsertTable returns:
// INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES
// ('it','mike','ssn_ca','Security Social Number, State of California','ssn,ca,california,wellfare',3,3),
func (o *Table) GetInsertTable() (string, error) {
	return fmt.Sprintf(`INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES(%s);`,
		o.String()), nil
}

// GetUpdateTable returns:
// UPDATE table_ SET descr='xxx',tags='yyyy',ncols=d,nrows=d
// WHERE owner='ooooo' AND name='nnnn'
func (o *Table) GetUpdateTable() (string, error) {
	return fmt.Sprintf(`UPDATE table_ SET descr='%s',tags='%s',ncols=%d,nrows=%d WHERE owner='%s' AND name='%s';`,
		o.Descr,
		o.Tags,
		o.NCols,
		o.NRows,
		o.Owner,
		o.Name), nil
}

// GetIncrementTable returns:
// UPDATE table_ SET nrows+=<param>
// WHERE owner='ooooo' AND name='nnnn'
func (o *Table) GetIncrementTable(nrows int64) (string, error) {
	if nrows <= 0 {
		return "", fmt.Errorf("table_.nrows update makes no sense because rows added were %d", nrows)
	}
	return fmt.Sprintf(`UPDATE table_ SET nrows=nrows+%d WHERE owner='%s' AND name='%s';`,
		nrows,
		o.Owner,
		o.Name), nil
}
