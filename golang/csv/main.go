package main

import (
	"bufio"
	"csv/taosUtil"
	"encoding/csv"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "hlhz",
}
var wg sync.WaitGroup

func main() {
	rootDir := "D:/csv/"
	files, _ := ioutil.ReadDir(rootDir)
	execFiles := make([]string, 0)
	for _, f := range files {
		if strings.Contains(f.Name(), ".csv") {
			s := strings.Split(f.Name(), "-")[1]
			if len(s) != 17 {
				fmt.Println(f.Name() + "文件开始时间无效")
				return
			}
			execFiles = append(execFiles, rootDir+f.Name())
		}
	}
	if len(execFiles) == 0 {
		fmt.Println("无文件可读")
		return
	}
	fmt.Printf("读取文件列表：%v\n", execFiles)
	for _, v := range execFiles {
		sinceStart := time.Now()
		stat, err := os.Stat(v)
		checkErr(err, "file stat error")
		start := strings.Split(stat.Name(), "-")[1]
		start = strings.Split(start, ".")[0]
		startTime, err := strconv.ParseInt(start, 10, 64)
		checkErr(err, "开始时间格式不正确")
		csvFile, err := os.Open(v)
		checkErr(err, "file open error")
		reader := csv.NewReader(bufio.NewReaderSize(csvFile, 1024*1024*100))
		lines, err := reader.ReadAll()
		checkErr(err, "read error")
		indexRecode := reader.FieldsPerRecord
		csvFile.Close()
		if indexRecode == 3 {
			writeTaosVo(lines, startTime)
		} else if indexRecode == 10 {
			writeTaosCur(lines, startTime)
		}
		fmt.Printf("总耗时：%v\n", time.Since(sinceStart))
		delFiles([]string{v})
		fmt.Println("删除文件:" + v)
	}

}

func writeTaosCur(lines [][]string, startTime int64) {
	startTime = startTime * 1000 //to us
	var result = make([]taosUtil.SubTableValue, 0)
	oneInsertSize := 1000
	rowValues := make([]taosUtil.RowValue, 0)
	for i := range lines {
		if i == 0 {
			continue
		}
		line := lines[i]
		ts, v0, v1, v2, v3, v4, v5, v6, v7, v8 := line[0], line[1], line[2], line[3], line[4], line[5], line[6], line[7], line[8], line[9]
		tsFloat, _ := strconv.ParseFloat(ts, 64)
		tsFloat = tsFloat * 1e6
		exeTs := int64(tsFloat) + startTime
		fieldValues := []taosUtil.FieldValue{
			{
				Name:  "ts",
				Value: exeTs,
			}, {
				Name:  "value0",
				Value: v0,
			}, {
				Name:  "value1",
				Value: v1,
			}, {
				Name:  "value2",
				Value: v2,
			}, {
				Name:  "value3",
				Value: v3,
			}, {
				Name:  "value4",
				Value: v4,
			}, {
				Name:  "value5",
				Value: v5,
			}, {
				Name:  "value6",
				Value: v6,
			}, {
				Name:  "value7",
				Value: v7,
			}, {
				Name:  "value8",
				Value: v8,
			},
		}
		rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
		if len(rowValues) >= oneInsertSize {
			var subTableValue taosUtil.SubTableValue
			subTableValue.Name = "c0"
			subTableValue.SuperTable = "cur"
			tags := []taosUtil.TagValue{
				{
					Name:  "groupId",
					Value: "0",
				},
			}
			subTableValue.Tags = tags
			subTableValue.Values = rowValues
			result = append(result, subTableValue)
			rowValues = make([]taosUtil.RowValue, 0)
		}
	}
	fmt.Println(len(result))
	groupSize := len(result)/2 + 1 // 每个线程执行40个对象
	startInsert(groupSize, result)
}

func writeTaosVo(lines [][]string, startTime int64) {
	startTime = startTime * 1000 //to us
	var result = make([]taosUtil.SubTableValue, 0)
	oneInsertSize := 20000
	rowValues := make([]taosUtil.RowValue, 0)
	for i := range lines {
		if i == 0 {
			continue
		}
		line := lines[i]
		ts := line[0]
		v0 := line[1]
		v1 := line[2]
		tsFloat, _ := strconv.ParseFloat(ts, 64)
		tsFloat = tsFloat * 1e6
		exeTs := int64(tsFloat) + startTime
		fieldValues := []taosUtil.FieldValue{
			{
				Name:  "ts",
				Value: exeTs,
			}, {
				Name:  "value0",
				Value: v0,
			}, {
				Name:  "value1",
				Value: v1,
			},
		}
		rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
		if len(rowValues) >= oneInsertSize {
			var subTableValue taosUtil.SubTableValue
			subTableValue.Name = "v0"
			subTableValue.SuperTable = "vo"
			tags := []taosUtil.TagValue{
				{
					Name:  "groupId",
					Value: "0",
				},
			}
			subTableValue.Tags = tags
			subTableValue.Values = rowValues
			result = append(result, subTableValue)
			rowValues = make([]taosUtil.RowValue, 0)
		}
	}
	groupSize := len(result)/2 + 1 // 每个线程执行40个对象
	fmt.Println(len(result))
	startInsert(groupSize, result)
}
func insert(array []taosUtil.SubTableValue) {
	defer wg.Done()
	for i := range array {
		_, err := taosUtil.InsertAutoCreateTable(array[i : i+1])
		checkErr(err, "taos insert error")
		//time.Sleep(time.Millisecond * 200)
	}
}
func startInsert(groupSize int, result []taosUtil.SubTableValue) {
	max := len(result)
	var quantity int
	if max%groupSize == 0 {
		quantity = max / groupSize
	} else {
		quantity = (max / groupSize) + 1
	}
	var start, end, i int
	for i = 1; i <= quantity; i++ {
		end = i * groupSize
		if i != quantity {
			wg.Add(1)
			go insert(result[start:end])
		} else {
			wg.Add(1)
			go insert(result[start:])
		}
		start = i * groupSize
	}
	wg.Wait()
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
