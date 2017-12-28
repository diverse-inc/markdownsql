package sqldef

import (
	"fmt"
	"regexp"
	"strings"
)

type SQLDefinition struct {
	DatabaseName string
	Tables       []TableDefinition
}

func (sd SQLDefinition) ToSQLStmt() string {
	sql := ``
	for _, tbl := range sd.Tables {
		// sql = sql + fmt.Sprintf("drop table if exists `%s`;\n", tbl.Name)
		sql = sql + fmt.Sprintf("create table `%s` (\n", tbl.Name)

		length := len(tbl.Columns)
		for i, column := range tbl.Columns {
			comma := `,`
			if !tbl.HasIndexes() && i == length-1 {
				comma = ``
			}

			sql = sql + fmt.Sprintf("  %s%s\n", column.ToSQLStmt(), comma)
		}

		length = len(tbl.Indexes)
		for i, index := range tbl.Indexes {
			comma := `,`
			if i == length-1 {
				comma = ``
			}

			sql = sql + fmt.Sprintf("  %s%s\n", index.ToSQLStmt(tbl.Name), comma)
		}

		sql = sql + fmt.Sprintln(`) /*! engine=innodb default charset=utf8mb4 row_format=dynamic */ ;`)
		sql = sql + fmt.Sprintln(``)
	}

	return sql
}

type TableDefinition struct {
	Name    string
	Columns []ColumnDefinition
	Indexes []IndexDefinition
}

func (tbl TableDefinition) HasIndexes() bool {
	return 0 < len(tbl.Indexes)
}

type ColumnDefinition struct {
	Name    string
	Type    string
	Null    bool
	Default string
	Key     string
	Comment string
}

func (cd ColumnDefinition) ToSQLStmt() string {
	null := `not null`
	if cd.Null {
		null = `null`
	}

	dflt := ``
	if 0 < len(cd.Default) {
		dflt = fmt.Sprintf("default %s", cd.Default)
	}

	key := ``
	if 0 < len(cd.Key) {
		key = cd.Key
	}

	comment := ``
	if 0 < len(cd.Comment) {
		comment = fmt.Sprintf("comment '%s'", cd.Comment)
	}

	s := fmt.Sprintf("`%s` %s %s %s %s %s", cd.Name, cd.Type, null, dflt, key, comment)
	s = strings.Trim(s, ` `)
	rgx := regexp.MustCompile(" +")
	return rgx.ReplaceAllString(s, ` `)
}

type IndexDefinition struct {
	Name     string
	Columns  []string
	IsUnique bool
}

func (idx IndexDefinition) ToSQLStmt(name string) string {
	columns := strings.Join(idx.Columns, "`,`")
	columns = strings.Replace(columns, " ", "", -1)
	if idx.IsUnique {
		return fmt.Sprintf("unique key %s_%s(`%s`)", idx.Name, name, columns)
	}
	return fmt.Sprintf("key %s_%s(`%s`)", idx.Name, name, columns)
}
