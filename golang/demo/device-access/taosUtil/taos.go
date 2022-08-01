package taosUtil

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v2/taosSql"
	"strconv"
	"strings"
)

type TaosInfo struct {
	HostName   string
	ServerPort int
	User       string
	Password   string
	DbName     string
}

type SubTableValue struct {
	SuperTable string
	Name       string
	Tags       []TagValue
	Values     []RowValue
}

type TagValue struct {
	Name  string
	Value interface{}
}

type RowValue struct {
	Fields []FieldValue
}

type FieldValue struct {
	Name  string
	Value interface{}
}

var TaosDb *sql.DB

func Connection(info TaosInfo) (*sql.DB, error) {
	url := "root:taosdata@/tcp(" + info.HostName + ":" + strconv.Itoa(info.ServerPort) + ")/" + info.DbName
	url += "?charset=utf-8&locale=en_US.UTF-8&timezone=UTC-8&maxSQLLength=1048576"
	var error error
	TaosDb, error = sql.Open("taosSql", url)
	return TaosDb, error
}

func Close() {
	TaosDb.Close()
}

func InsertAutoCreateTable(subTableValue []SubTableValue) (int64, error) {
	sqlStr := insertMultiTableMultiValuesUsingSuperTable(subTableValue)
	res, err := TaosDb.Exec(sqlStr)
	if err != nil {
		return -1, err
	}
	return res.RowsAffected()
}

func insertMultiTableMultiValuesUsingSuperTable(subTableValue []SubTableValue) string {
	head := "insert into "
	tableArray := make([]string, len(subTableValue))
	for i, tableValue := range subTableValue {
		var buffer bytes.Buffer
		buffer.WriteString(tableValue.Name)
		buffer.WriteString(" using ")
		buffer.WriteString(tableValue.SuperTable)
		buffer.WriteString(" tags ")
		buffer.WriteString(tagValues(tableValue.Tags))
		buffer.WriteString(" values ")
		buffer.WriteString(rowValues(tableValue.Values))
		tableArray[i] = buffer.String()
	}
	return head + strings.Join(tableArray, " ")
}

func rowValues(values []RowValue) string {
	tableArray := make([]string, len(values))
	for i, value := range values {
		tableArray[i] = fieldValues(value.Fields)
	}
	return strings.Join(tableArray, ",")
}

func fieldValues(fields []FieldValue) string {
	array := make([]string, len(fields))
	for i := range array {
		if i == 0 {
			array[i] = fmt.Sprintf("%v", fields[i].Value)
		} else {
			array[i] = fmt.Sprintf("'%v'", fields[i].Value)
		}
	}
	return "(" + strings.Join(array, ",") + ")"
}

func tagValues(tags []TagValue) string {
	tableArray := make([]string, len(tags))
	for i, tag := range tags {
		tableArray[i] = fmt.Sprintf("'%v'", tag.Value)
	}
	return "(" + strings.Join(tableArray, ",") + ")"
}

func ExecuteQuery(sqlStr string) ([]map[string]string, error) {
	rows, err := TaosDb.Query(sqlStr)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	values := make([]sql.RawBytes, len(columns))
	scanArgs := make([]interface{}, len(values))
	ret := make([]map[string]string, 0)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			break
		}
		var value string
		vmap := make(map[string]string, len(scanArgs))
		for i, col := range values {
			if col == nil {
				value = "" // or NULL
			} else {
				value = string(col)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	return ret, err
}

func UseDatabase(database string) error {
	stmt, err := TaosDb.Prepare("use " + database)
	if err != nil {
		return err
	}
	res, err := stmt.Exec()
	if err != nil {
		return err
	}
	_, err = res.RowsAffected()
	return err
}
