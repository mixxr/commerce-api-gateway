package models

import (
	"bytes"
	"fmt"
	"strings"
)

const COLNAMES_SUFFIX string = "colnames"

type TableColnames struct {
	table_ *Table
	Lang   string   `json:"lang"`
	Header []string `json:"header"`
}

func (o *TableColnames) String() string {
	return fmt.Sprintf("%s;%s;", o.Lang, strings.Join(o.Header, ","))
}

func NewColnames(o *Table, header []string) *TableColnames {
	o.NCols = len(header)
	newObj := new(TableColnames)
	newObj.Lang = o.DefLang
	newObj.table_ = o
	newObj.Header = header
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

	fmt.Fprintf(&buffer, "CREATE TABLE %s_%s_%s (", o.table_.Owner, o.table_.Name, COLNAMES_SUFFIX)

	buffer.WriteString("id BIGINT NOT NULL AUTO_INCREMENT,")
	buffer.WriteString("created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,")
	buffer.WriteString("lang CHAR(2),")

	for k, _ := range o.Header {
		fmt.Fprintf(&buffer, "colname%d VARCHAR(32) NOT NULL,", k)
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

	fmt.Fprintf(&buffer, "INSERT INTO %s_%s_%s (lang", o.table_.Owner, o.table_.Name, COLNAMES_SUFFIX)
	for k, _ := range o.Header {
		fmt.Fprintf(&buffer, ",colname%d", k)
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
	return fmt.Sprintf("DELETE %s_%s_%s WHERE lang='%s';",
		o.table_.Owner,
		o.table_.Name,
		COLNAMES_SUFFIX,
		o.Lang), nil

}
