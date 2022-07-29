package mysqlUtil

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
)

type MysqlInfo struct {
	UserName string
	Password string
	Ip       string
	Port     int
	Db       string
}

var MyDB *sql.DB

func Connection(info MysqlInfo) (*sql.DB, error) {
	path := strings.Join([]string{info.UserName, ":", info.Password, "@tcp(", info.Ip, ":", strconv.Itoa(info.Port), ")/", info.Db, "?charset=utf8"}, "")
	var err error
	MyDB, err = sql.Open("mysql", path)
	if err != nil {
		return nil, err
	}
	MyDB.SetConnMaxLifetime(100)
	MyDB.SetMaxOpenConns(100)
	MyDB.SetMaxIdleConns(10)
	error := MyDB.Ping()
	return MyDB, error
}

//Add new record to table
func Add(insertSql string, args ...interface{}) (int64, error) {
	stmt, err := MyDB.Prepare(insertSql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.LastInsertId()
}

//Exe commands
//del or update
func Exe(exeSql string, args ...interface{}) (int64, error) {
	stmt, err := MyDB.Prepare(exeSql)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}

//GetFirst
//get one record from query result
func GetFirst(query string, args ...interface{}) (map[string]string, error) {
	if !strings.Contains(strings.ToUpper(query), "LIMIT") {
		query += " LIMIT 1"
	}
	stmt, err := MyDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
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
	ret := make(map[string]string, len(scanArgs))

	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			break
		}
		var value string
		for i, col := range values {
			if col == nil {
				value = "" //or NULL
			} else {
				value = string(col)
			}
			ret[columns[i]] = value
		}
		break //get the first row only
	}
	return ret, err
}

//GetAll
//all records from query result
func GetAll(query string, args ...interface{}) ([]map[string]string, error) {
	stmt, err := MyDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	rows, err := stmt.Query(args...)
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
