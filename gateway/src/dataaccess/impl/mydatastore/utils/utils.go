package utils

import (
	"bytes"
	"fmt"
	"main/dataaccess/models"
	"strings"
)

// =============== table

// GetIncrementTable returns:
// UPDATE table_ SET nrows+=<param>
// WHERE owner='ooooo' AND name='nnnn'
func GetIncrementTable(o *models.Table, nrows int64) (string, error) {
	if nrows == 0 {
		return "", fmt.Errorf("table_.nrows update makes no sense because rows added are 0")
	}
	operation := "+"
	if nrows < 0 {
		operation = "-"
		nrows = -nrows
	}
	return fmt.Sprintf(`UPDATE table_ SET nrows=nrows%s%d WHERE owner='%s' AND name='%s';`,
		operation,
		nrows,
		o.Owner,
		o.Name), nil
}

// GetSelectNRows returns:
// SELECT NRows as tot from table_ WHERE owner='%s' AND name='%s'
func GetSelectNRows(tin *models.Table) (string, error) {
	return fmt.Sprintf(`SELECT NRows as tot from table_ WHERE owner='%s' AND name='%s';`,
		tin.Owner,
		tin.Name), nil
}

// GetInsertTable returns:
// INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES
// ('it','mike','ssn_ca','Security Social Number, State of California','ssn,ca,california,wellfare',3,3),
func GetInsertTable(o *models.Table) (string, error) {
	return fmt.Sprintf(`INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows,status) VALUES(%s);`,
		o.String()), nil
}

// GetUpdateTable returns:
// UPDATE table_ SET descr='xxx',tags='yyyy',ncols=d,nrows=d, status
// WHERE owner='ooooo' AND name='nnnn'
func GetUpdateTable(o *models.Table) (string, error) {
	return fmt.Sprintf(`UPDATE table_ SET descr='%s',tags='%s',ncols=%d,nrows=%d,status=%d WHERE owner='%s' AND name='%s';`,
		o.Descr,
		o.Tags,
		o.NCols,
		o.NRows,
		o.Status,
		o.Owner,
		o.Name), nil
}

// GetSelectSearchTable returns:
// SELECT Id, Descr, Tags, DefLang, NCols, NRows FROM table_ WHERE Name like '%s<name>%s' and Owner like '%s<ownwer>%s' and ...
func GetSelectSearchTable(tin *models.Table) (string, error) {
	if tin.IsEmpty() {
		return "", fmt.Errorf("service search makes no sense because all parameters empty")
	}

	var where []string

	if tin.Name != "" && tin.Name != "-" {
		tin.Name = strings.Replace(tin.Name, "-", "%", 1)
		tin.Name = strings.ReplaceAll(tin.Name, "_", "\\_")
		where = append(where, fmt.Sprintf("name like '%s'", tin.Name))
	}
	if tin.Owner != "" && tin.Owner != "-" {
		tin.Owner = strings.Replace(tin.Owner, "-", "%", 1)
		tin.Owner = strings.ReplaceAll(tin.Owner, "_", "\\_")
		where = append(where, fmt.Sprintf("owner like '%s'", tin.Owner))
	}
	if tin.Descr != "" && tin.Descr != "-" {
		tin.Descr = strings.ReplaceAll(tin.Descr, "-", "%")
		tin.Descr = strings.ReplaceAll(tin.Descr, "_", "\\_")
		tin.Descr = strings.ReplaceAll(tin.Descr, " ", "_") // replace spaces with wildcard _
		where = append(where, fmt.Sprintf("descr like '%s'", tin.Descr))
	}
	where = append(where, fmt.Sprintf("status>=%d", tin.Status))

	return fmt.Sprintf(`SELECT Id, Owner, Name, Descr, Tags, DefLang, NCols, NRows FROM table_ WHERE %s;`,
		strings.Join(where, " AND ")), nil
}

func GetSelectTable(name string, owner string, status int) (string, error) {
	if name == "" || owner == "" || !models.IsValidStatus(status) {
		return "", fmt.Errorf("cannot select with wrong parameters: %s %s %d", name, owner, status)
	}
	return fmt.Sprintf(`SELECT Id, Descr, Tags, DefLang, NCols, NRows, Status FROM table_ WHERE Name='%s' and Owner='%s' and status>=%d;`,
		name,
		owner,
		status), nil
}

// GetUpdateNCols returns:
// UPDATE table_ SET ncols=Table.NCols
// WHERE owner='ooooo' AND name='nnnn'
func GetUpdateNCols(o *models.Table, ncols int) (string, error) {
	if ncols <= 0 {
		return "", fmt.Errorf("table_.ncols update makes no sense because ncols is %d", ncols)
	}
	return fmt.Sprintf(`UPDATE table_ SET ncols=%d WHERE owner='%s' AND name='%s';`,
		ncols,
		o.Owner,
		o.Name), nil
}

// =========== TableColnames

const COLNAME_SUFFIX string = "colname"

func GetTableColnamesName(o *models.Table) string {
	if o == nil || o.Owner == "" || o.Name == "" {
		// TODO: exception propagation ?
		return "" //, fmt.Errorf("table name not defined, pls set table_: %s", o.Parent())
	}
	// status makes the table name
	return fmt.Sprintf("%s_%s_%ss_%d", o.Owner, o.Name, COLNAME_SUFFIX, o.Status)
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
func GetCreateTableColnames(o *models.TableColnames) (string, error) {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "CREATE TABLE IF NOT EXISTS %s (", GetTableColnamesName(o.Parent()))

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
func GetInsertTableColnames(o *models.TableColnames) (string, error) {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "INSERT INTO %s (lang", GetTableColnamesName(o.Parent()))
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

// // GetDeleteTable returns SQL instruction: DELETE <owner>_<name>_colnames ... WHERE lang=
// func GetDeleteTableColnames(o *models.TableColnames) (string, error) {
// 	return fmt.Sprintf("DELETE FROM %s WHERE lang='%s';",
// 		GetTableColnamesName(o.Parent()),
// 		o.Lang), nil

// }

func GetSelectTableColnames(o *models.TableColnames) (string, error) {
	if o.Parent() == nil {
		return "", fmt.Errorf("cannot define a SELECT with table_=nil")
	}
	if o.Lang == "" {
		return "", fmt.Errorf("colnames select makes no sense without a lang: %s", o.Lang)
	}
	if o.Parent().NCols <= 0 {
		return "", fmt.Errorf("cannot define a SELECT with table_NCols=%d", o.Parent().NCols)
	}

	var buffer bytes.Buffer

	buffer.WriteString("SELECT lang")
	for i := 0; i < o.Parent().NCols; i++ {
		fmt.Fprintf(&buffer, ",%s%d", COLNAME_SUFFIX, i)
	}
	fmt.Fprintf(&buffer, " FROM %s WHERE lang='%s';",
		GetTableColnamesName(o.Parent()),
		o.Lang)

	return buffer.String(), nil
}

// =========== TableValues

const VALUE_SUFFIX string = "value"

func GetTableValuesName(o *models.Table) string {
	if o == nil || o.Owner == "" || o.Name == "" {
		// TODO: exception propagation ?
		return "" //, fmt.Errorf("table name not defined, pls set table_: %s", o.Parent())
	}
	return fmt.Sprintf("%s_%s_%ss_%d", o.Owner, o.Name, VALUE_SUFFIX, o.Status)
}

// returns SQL instruction: CREATE TABLE <owner>_<name>_values ...
// example
// CREATE TABLE mike_ssn_ca_values (
// 	id BIGINT NOT NULL AUTO_INCREMENT,
// 	created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
// 	colvalue1 VARCHAR(256) NOT NULL,
// 	colvalue2 VARCHAR(256) NOT NULL,
// 	colvalue3 VARCHAR(256) NOT NULL,
// 	PRIMARY KEY ( id )
//  );
func GetCreateTableValues(o *models.TableValues) (string, error) {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "CREATE TABLE IF NOT EXISTS %s (", GetTableValuesName(o.Parent()))

	buffer.WriteString("id BIGINT NOT NULL AUTO_INCREMENT,")
	buffer.WriteString("created TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,")

	for k := 0; k < len(o.Rows[0]); k++ {
		fmt.Fprintf(&buffer, "%s%d VARCHAR(256) NOT NULL,", VALUE_SUFFIX, k)
	}

	buffer.WriteString("PRIMARY KEY ( id ));")

	return buffer.String(), nil
}

// returns SQL instruction: INSERT INTO  <owner>_<name>_values ...
// Example:
// INSERT INTO mike_ssn_ca (colvalue1,colvalue2,colvalue3) VALUES
// ('mike','douglàs','3897428934EWREW'),
// ('äbel','òmar ópël','3897428934EWREW'),
// ('anthony','di martino','234234FSAFSADF');

func GetInsertTableValues(o *models.TableValues) (string, error) {
	var buffer bytes.Buffer

	fmt.Fprintf(&buffer, "INSERT INTO %s (", GetTableValuesName(o.Parent()))
	i, j := 0, 0
	for ; i < len(o.Rows[0])-1; i++ {
		fmt.Fprintf(&buffer, "%s%d,", VALUE_SUFFIX, i)
	}
	fmt.Fprintf(&buffer, "%s%d) VALUES ", VALUE_SUFFIX, i)
	i = 0
	for ; i < len(o.Rows)-1; i++ {
		buffer.WriteString("(")
		j = 0
		for ; j < len(o.Rows[i])-1; j++ {
			fmt.Fprintf(&buffer, "'%s',", o.Rows[i][j])
		}
		fmt.Fprintf(&buffer, "'%s'),", o.Rows[i][j])
	}
	buffer.WriteString("(")
	j = 0
	for ; j < len(o.Rows[i])-1; j++ {
		fmt.Fprintf(&buffer, "'%s',", o.Rows[i][j])
	}
	fmt.Fprintf(&buffer, "'%s');", o.Rows[i][j])

	return buffer.String(), nil
}

func GetSelectTableValues(o *models.TableValues) (string, error) {
	if o.Parent() == nil {
		return "", fmt.Errorf("cannot define a SELECT with table_=nil")
	}
	if o.Start < 0 {
		return "", fmt.Errorf("cannot define a SELECT with start<0")
	}
	if o.Count <= 0 {
		return "", fmt.Errorf("cannot define a SELECT with count less than 0")
	}
	if o.Parent().NCols <= 0 {
		return "", fmt.Errorf("cannot define a SELECT with NCols<=0")
	}

	var buffer bytes.Buffer

	buffer.WriteString("SELECT ")
	for i := 0; i < o.Parent().NCols-1; i++ {
		fmt.Fprintf(&buffer, "%s%d,", VALUE_SUFFIX, i)
	}
	fmt.Fprintf(&buffer, "%s%d", VALUE_SUFFIX, o.Parent().NCols-1)
	fmt.Fprintf(&buffer, " FROM %s LIMIT %d, %d;",
		GetTableValuesName(o.Parent()),
		o.Start,
		o.Count)

	return buffer.String(), nil
}

// DELETE

// GetDropTables DROP both <owner>_<name>_colnames and <owner>_<name>_values in NOWAIT
func GetDropTables(o *models.Table) (string, error) {
	// to be sure to remove any table
	o.Status = models.StatusDraft
	tableNames := []string{
		GetTableColnamesName(o),
		GetTableValuesName(o),
	}
	o.Status = models.StatusEnabled
	return fmt.Sprintf("DROP TABLE IF EXISTS %s, %s, %s, %s NOWAIT;",
		tableNames[0],
		tableNames[1],
		GetTableColnamesName(o),
		GetTableValuesName(o)), nil
}

func GetDeleteTable(o *models.Table) (string, error) {
	return fmt.Sprintf("DELETE QUICK FROM table_ WHERE owner='%s' AND name='%s';",
		o.Owner,
		o.Name), nil
}

func GetDeleteTableColnames(o *models.Table, langs []string) (string, error) {
	var buffer bytes.Buffer
	if len(langs) > 0 {
		for i, lang := range langs {
			langs[i] = fmt.Sprintf("'%s'", lang)
		}
		fmt.Fprintf(&buffer, " WHERE lang IN (%s);", strings.Join(langs, ","))
	}
	return fmt.Sprintf("DELETE QUICK FROM %s %s;",
		GetTableColnamesName(o),
		buffer.String()), nil
}

func GetDeleteTableValues(o *models.Table, count int64) (string, error) {
	var buffer bytes.Buffer
	if count != 0 {
		if count < 0 {
			count = -count
			buffer.WriteString(" ORDER BY id DESC")
		}
		fmt.Fprintf(&buffer, " LIMIT %d;", count)
	}
	return fmt.Sprintf("DELETE FROM %s %s;",
		GetTableValuesName(o),
		buffer.String()), nil
}

// GetRenameTables returns RENAME TABLE %s TO %s, %s TO %s;
func GetRenameTables(o *models.Table, newStatus int) (string, error) {
	oldTableCN := GetTableColnamesName(o)
	oldTableV := GetTableValuesName(o)
	o.Status = newStatus
	return fmt.Sprintf("RENAME TABLE %s TO %s, %s TO %s;",
		oldTableCN,
		GetTableColnamesName(o),
		oldTableV,
		GetTableValuesName(o)), nil
}
