package utils

import (
	"bytes"
	"dataaccess/models"
	"fmt"
)

// =============== table

// GetIncrementTable returns:
// UPDATE table_ SET nrows+=<param>
// WHERE owner='ooooo' AND name='nnnn'
func GetIncrementTable(o *models.Table, nrows int64) (string, error) {
	if nrows <= 0 {
		return "", fmt.Errorf("table_.nrows update makes no sense because rows added were %d", nrows)
	}
	return fmt.Sprintf(`UPDATE table_ SET nrows=nrows+%d WHERE owner='%s' AND name='%s';`,
		nrows,
		o.Owner,
		o.Name), nil
}

// GetInsertTable returns:
// INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES
// ('it','mike','ssn_ca','Security Social Number, State of California','ssn,ca,california,wellfare',3,3),
func GetInsertTable(o *models.Table) (string, error) {
	return fmt.Sprintf(`INSERT INTO table_ (deflang,owner,name,descr,tags,ncols,nrows) VALUES(%s);`,
		o.String()), nil
}

// GetUpdateTable returns:
// UPDATE table_ SET descr='xxx',tags='yyyy',ncols=d,nrows=d
// WHERE owner='ooooo' AND name='nnnn'
func GetUpdateTable(o *models.Table) (string, error) {
	return fmt.Sprintf(`UPDATE table_ SET descr='%s',tags='%s',ncols=%d,nrows=%d WHERE owner='%s' AND name='%s';`,
		o.Descr,
		o.Tags,
		o.NCols,
		o.NRows,
		o.Owner,
		o.Name), nil
}

// GetSelectTable returns:
// SELECT Id, Descr, Tags, DefLang, NCols, NRows FROM table_ WHERE Name='%s' and Owner='%s'
func GetSelectTable(name string, owner string) (string, error) {
	if name == "" || owner == "" {
		return "", fmt.Errorf("table_ select makes no sense because name %s and/or owner %s were empty", name, owner)
	}
	return fmt.Sprintf(`SELECT Id, Descr, Tags, DefLang, NCols, NRows FROM table_ WHERE Name='%s' and Owner='%s';`,
		name,
		owner), nil
}

// GetUpdateNCols returns:
// UPDATE table_ SET ncols=Table.NCols
// WHERE owner='ooooo' AND name='nnnn'
func GetUpdateNCols(o *models.Table) (string, error) {
	ncols := o.NCols
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

func GetTableColnamesName(o *models.TableColnames) string {
	if o.Parent() == nil || o.Parent().Owner == "" || o.Parent().Name == "" {
		// TODO: exception propagation ?
		return "" //, fmt.Errorf("table name not defined, pls set table_: %s", o.Parent())
	}
	return fmt.Sprintf("%s_%s_%ss", o.Parent().Owner, o.Parent().Name, COLNAME_SUFFIX)
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

	fmt.Fprintf(&buffer, "CREATE TABLE %s (", GetTableColnamesName(o))

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

	fmt.Fprintf(&buffer, "INSERT INTO %s (lang", GetTableColnamesName(o))
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
func GetDeleteTableColnames(o *models.TableColnames) (string, error) {
	return fmt.Sprintf("DELETE %s WHERE lang='%s';",
		GetTableColnamesName(o),
		o.Lang), nil

}

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
		GetTableColnamesName(o),
		o.Lang)

	return buffer.String(), nil
}

// =========== TableValues

const VALUE_SUFFIX string = "value"

func GetTableValuesName(o *models.TableValues) string {
	if o.Parent() == nil || o.Parent().Owner == "" || o.Parent().Name == "" {
		// TODO: exception propagation ?
		return "" //, fmt.Errorf("table name not defined, pls set table_: %s", o.Parent())
	}
	return fmt.Sprintf("%s_%s_%ss", o.Parent().Owner, o.Parent().Name, VALUE_SUFFIX)
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

	fmt.Fprintf(&buffer, "CREATE TABLE %s (", GetTableValuesName(o))

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

	fmt.Fprintf(&buffer, "INSERT INTO %s (", GetTableValuesName(o))
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
	if o.Start <= 0 {
		return "", fmt.Errorf("cannot define a SELECT with start<=0")
	}
	if o.Count <= 0 {
		return "", fmt.Errorf("cannot define a SELECT with count<=0")
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
	fmt.Fprintf(&buffer, " FROM %s WHERE id>=%d AND id<%d", //LIMIT %d;",
		GetTableValuesName(o),
		o.Start,
		o.Count+o.Start)

	return buffer.String(), nil
}
