package main

import (
	"excel/taosUtil"
	"fmt"
	"github.com/xuri/excelize/v2"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "hlhz",
}
var layout = "2006/01/02 15:04:05"
var fieldDir = map[string]string{
	"采集时间":     "ts",
	"制丝车间环境温度": "zscjhjwd",
	"制丝车间环境湿度": "zscjhjsd",
	"送料高度":     "slgd",
	"上刀门压力":    "sdmyl",
	"储叶间环境温度":  "cyjhjwd",
	"储叶间环境湿度":  "cyjhjsd",
	"砂轮使用时长":   "slsysc",
	"砂轮磨刀时间间隔": "slmdsjjg",
}

func main() {
	rootDir := "D:/excel/"
	files, _ := ioutil.ReadDir(rootDir)
	execFiles := make([]string, 0)
	for _, f := range files {
		if strings.Contains(f.Name(), ".xlsx") {
			execFiles = append(execFiles, rootDir+f.Name())
		}
	}
	if len(execFiles) == 0 {
		fmt.Println("无文件可读")
		return
	}
	fmt.Printf("读取文件列表：%v\n", execFiles)
	for _, v := range execFiles {
		f, err := excelize.OpenFile(v)
		checkErr(err, "打开文件失败:"+v)
		sheetList := f.GetSheetList()
		for _, sheet := range sheetList {
			fieldNames := make([]string, 0)
			rows, err := f.Rows(sheet)
			checkErr(err, "获取sheet里面的数据失败:"+v)
			rowIndex := 0
			for rows.Next() {
				row, err := rows.Columns()
				checkErr(err, "获取行数据失败:"+v)
				var ts int64
				subTables := make([]taosUtil.SubTableValue, 0)
				for index, colCell := range row {
					if rowIndex == 0 {
						fieldNames = append(fieldNames, fieldDir[colCell])
					} else {
						if index == 0 {
							excelTime := colCell
							excelTimeLocal, err := time.ParseInLocation(layout, excelTime, time.UTC)
							checkErr(err, "转换时间出错")
							ts = excelTimeLocal.UnixNano() / 1e3
						} else {
							tableName := fieldNames[index]
							subTable := getSubTableValue(ts, tableName, colCell)
							subTables = append(subTables, subTable)
						}
					}
				}
				rowIndex++
				if len(subTables) > 0 {
					_, err = taosUtil.InsertAutoCreateTable(subTables)
					checkErr(err, "insert taos error:"+v+","+sheet)
				}
			}
			err = rows.Close()
			checkErr(err, "row 关闭失败.")
		}
		err = f.Close()
		checkErr(err, "关闭excel失败:"+v)
		/*delFiles([]string{v})
		fmt.Println("删除文件:" + v)*/
	}
}

func getSubTableValue(ts int64, tableName string, value string) taosUtil.SubTableValue {
	var subTableValue taosUtil.SubTableValue
	subTableValue.Name = tableName
	subTableValue.SuperTable = "excel"
	tags := []taosUtil.TagValue{
		{
			Name:  "fname",
			Value: tableName,
		},
	}

	subTableValue.Tags = tags
	subTableValue.Values = []taosUtil.RowValue{{
		Fields: []taosUtil.FieldValue{
			{
				Name:  "ts",
				Value: ts,
			}, {
				Name:  "value",
				Value: value,
			},
		},
	}}
	return subTableValue
}
func delFiles(files []string) error {
	if len(files) == 0 {
		fmt.Println("需要删除的文件为空")
		return nil
	}
	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Println(file + "文件删除失败:" + err.Error())
			continue
		}
	}
	return nil
}
func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
func init() {
	taosUtil.Connection(taosInfo)
}
