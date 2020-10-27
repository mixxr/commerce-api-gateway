package models

import (
	"bytes"
	"fmt"
	"strings"
)

const COLNAME_SUFFIX string = "colname"

func (o *TableColnames) GetTableName() string {
	if o.table_ == nil || o.table_.Owner == "" || o.table_.Name == "" {
		// TODO: exception propagation ?
		return "" //, fmt.Errorf("table name not defined, pls set table_: %s", o.table_)
	}
	return fmt.Sprintf("%s_%s_%ss", o.table_.Owner, o.table_.Name, COLNAME_SUFFIX)
}

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

// returns SQL instruction: CREATE TABLE <owner>_<name>_colnames ...
// example
// CREATE TABLE mike_ssn_ca_colnames (
// 	id BIGINT NOT NULL AUTO_INCREMENT,
// 	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 	lang CHAR(2),
// 	colname1 VARCHAR(32) NOT NULL,
// 	colname2 VARCHAR(32) NOT NULL,
// 	colname3 VARCHAR(32) NOT NULL,
// 	PRIMARY KEY ( id )
//  );
func (o *TableColnames) GetCreateTable() (string, error) {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "CREATE TABLE %s (", o.GetTableName())

	buffer.WriteString("id BIGINT NOT NULL AUTO_INCREMENT,")
	buffer.WriteString("created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,")
	buffer.WriteString("lang CHAR(2),")

	for k, _ := range o.Header {
		fmt.Fprintf(&buffer, "%s%d VARCHAR(32) NOT NULL,", COLNAME_SUFFIX, k)
	}

	buffer.WriteString("PRIMARY KEY ( id ));")

	return buffer.String(), nil
}

// returns SQL instruction: INSERT INTO <owner>_<name>_colnames ...
// Example:
// INSERT INTO mike_ssn_ca_colnames (lang,colname1,colname2,colname3) VALUES
// ('it','nome','cognome','ssn'),
// ('en','name','surname','ssn');
func (o *TableColnames) GetInsertTable() (string, error) {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "INSERT INTO %s (lang", o.GetTableName())
	for k, _ := range o.Header {
		fmt.Fprintf(&buffer, ",%s%d", COLNAME_SUFFIX, k)
	}

	fmt.Fprintf(&buffer, ") VALUES ('%s'", o.Lang)

	for _, v := range o.Header {
		fmt.Fprintf(&buffer, ",'%s'", v)
	}

	buffer.WriteString(");")

	return buffer.String(), nil
}

// GetDeleteTable returns SQL instruction: DELETE <owner>_<name>_colnames ... WHERE lang=
func (o *TableColnames) GetDeleteTable() (string, error) {
	return fmt.Sprintf("DELETE %s WHERE lang='%s';",
		o.GetTableName(),
		o.Lang), nil

}

func (o *TableColnames) GetSelectTable() (string, error) {
	if o.table_ == nil {
		return "", fmt.Errorf("cannot define a SELECT with table_=nil")
	}
	if o.Lang == "" {
		return "", fmt.Errorf("colnames select makes no sense without a lang: %s", o.Lang)
	}
	if o.table_.NCols <= 0 {
		return "", fmt.Errorf("cannot define a SELECT with table_NCols=%d", o.table_.NCols)
	}

	var buffer bytes.Buffer

	buffer.WriteString("SELECT lang")
	for i := 0; i < o.table_.NCols; i++ {
		fmt.Fprintf(&buffer, ",%s%d", COLNAME_SUFFIX, i)
	}
	fmt.Fprintf(&buffer, " FROM %s WHERE lang='%s';",
		o.GetTableName(),
		o.Lang)

	return buffer.String(), nil
}
