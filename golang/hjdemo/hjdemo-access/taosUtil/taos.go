package taosUtil

import (
	"bytes"
	"database/sql"
	"fmt"
	_ "github.com/taosdata/driver-go/v3/taosSql"
	"math/rand"
	"strconv"
	"strings"
	"time"
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
var sourceDatabase string

func Connection(info TaosInfo) (*sql.DB, error) {
	url := "root:taosdata@/tcp(" + info.HostName + ":" + strconv.Itoa(info.ServerPort) + ")/" + info.DbName
	url += "?charset=utf-8&locale=en_US.UTF-8&timezone=UTC-8&maxSQLLength=1048576"
	var error error
	TaosDb, error = sql.Open("taosSql", url)
	sourceDatabase = info.DbName
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
		buffer.WriteString(" ")
		buffer.WriteString(rowFields(tableValue.Values))
		buffer.WriteString(" values ")
		buffer.WriteString(rowValues(tableValue.Values))
		tableArray[i] = buffer.String()
	}
	return head + strings.Join(tableArray, " ")
}

func rowFields(values []RowValue) string {
	row := values[0].Fields
	array := make([]string, len(row))
	for i, value := range row {
		array[i] = value.Name
	}
	return "(" + strings.Join(array, ",") + ")"
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

func ExecuteQuery(sqlStr string, database string) ([]map[string]string, error) {
	if database != "" {
		err := UseDatabase(database)
		if err != nil {
			return nil, err
		}
	}
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
			if columns[i] == "ts" {
				// 时间转时间戳 2022-08-19T17:56:57+08:00
				time, _ := time.ParseInLocation(time.RFC3339, value, time.Local)
				value = strconv.FormatInt(time.UnixNano()/1e6, 10)
			}
			vmap[columns[i]] = value
		}
		ret = append(ret, vmap)
	}
	if database != "" {
		UseDatabase(sourceDatabase)
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

func CreateTaosData() []SubTableValue {
	var subTableValues []SubTableValue
	startTime := "2022-07-10 10:00:00"

	startDate, _ := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)

	start := startDate.UnixNano() / 1e6
	// 100个电表，每块电表200个点位
	for i := 0; i < 1; i++ {
		var subTableValue SubTableValue
		is := strconv.Itoa(i)
		subTableValue.Name = "d00" + is
		subTableValue.SuperTable = "meters"
		tags := []TagValue{
			{
				Name:  "location",
				Value: "d00" + is,
			},
			{
				Name:  "groupId",
				Value: 1 + i,
			},
		}
		subTableValue.Tags = tags
		num := 100 // 一次插入多少数据
		rowValues := make([]RowValue, num)
		for j := 0; j < num; j++ {
			fieldNum := 200
			fieldValues := make([]FieldValue, 0)
			fieldValues = append(fieldValues, FieldValue{
				Name:  "ts",
				Value: start,
			})
			for a := 0; a < fieldNum; a++ {

				fieldName := "field" + strconv.Itoa(a)

				f := rand.Int() % 1000
				v := 200 + rand.Float32()

				var fieldValue FieldValue
				if a%2 == 0 {
					fieldValue = FieldValue{
						Name:  fieldName,
						Value: v,
					}
				} else {
					fieldValue = FieldValue{
						Name:  fieldName,
						Value: f,
					}
				}
				fieldValues = append(fieldValues, fieldValue)
			}
			rowValues[j] = RowValue{Fields: fieldValues}
			start += 1000
		}
		subTableValue.Values = rowValues
		subTableValues = append(subTableValues, subTableValue)
	}
	return subTableValues

}
