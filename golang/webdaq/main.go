package main

import (
	"encoding/binary"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/rs/xid"
	"math"
	"strconv"
	"strings"
	"time"
	"webdaq/redisUtil"
	"webdaq/taosUtil"
)

var redisInfo = redisUtil.RedisInfo{
	Ip:       "10.88.0.14",
	Port:     63799,
	Password: "xapp",
	Db:       0,
}
var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "webdaq",
}
var layout = "2006-01-02 15:04:05"

func init() {
	redisUtil.RedisInit(redisInfo)
	taosUtil.Connection(taosInfo)
}

func main() {
	fmt.Println("webdaq 采集启动")
	key := "daqjob"
	consumersGroup := "consumer-group"
	uniqueID := xid.New().String()
	rc := redisUtil.Redis{}
	err := rc.XGroupCreate(3, key, consumersGroup, "0")
	if err != nil {
		if err.Error() == "BUSYGROUP Consumer Group name already exists" {
			fmt.Println("group exists.")
		} else {
			checkErr(err, "创建消费组失败")
		}
	}
	for {
		entries, err := rc.XReadGroup(3, &redis.XReadGroupArgs{
			Group:    consumersGroup,
			Consumer: uniqueID,
			Streams:  []string{key, ">"},
			Count:    2,
			Block:    0,
			NoAck:    false,
		})
		checkErr(err, "read error")
		for i := 0; i < len(entries[0].Messages); i++ {
			messageID := entries[0].Messages[i].ID
			values := entries[0].Messages[i].Values
			err := dealData(values)
			checkErr(err, "数据处理异常")
			rc.XAck(3, key, consumersGroup, messageID)
		}
	}
}
func dealData(data map[string]interface{}) error {
	ts, err := strconv.ParseInt(data["Ts"].(string), 10, 64)
	if err != nil {
		return err
	}
	value := data["Value"].(string)
	/*float := val2float(value)
	channel := splitChannel(float, 2)*/
	channel := valSplitChannel(value, 2)
	/*fmt.Printf("ts:%v\n", time.Unix(ts/1000, 0).Format(layout))
	for k, v := range channel {
		fmt.Printf("channel:%v,data len : %v\n", k, len(v))
	}*/
	taosData := getTaosData(channel, ts)
	now := time.Now()
	_, err = taosUtil.InsertAutoCreateTable(taosData)
	fmt.Printf("taos 插入时间：%v\n", time.Since(now))
	checkErr(err, "taos 插入错误")
	return nil
}

func getTaosData(channel map[int][]byte, time int64) []taosUtil.SubTableValue {
	var result = make([]taosUtil.SubTableValue, 0)
	for k, v := range channel {
		var subTableValue taosUtil.SubTableValue
		itoa := strconv.Itoa(k)
		subTableValue.Name = "webdaq_" + itoa
		subTableValue.SuperTable = "webdaq"
		tags := []taosUtil.TagValue{
			{
				Name:  "channel",
				Value: k,
			},
		}
		subTableValue.Tags = tags
		rowValues := make([]taosUtil.RowValue, 1)
		values := val2float(v)
		fieldValues := []taosUtil.FieldValue{
			{
				Name:  "ts",
				Value: time,
			}, {
				Name:  "value0",
				Value: values[0],
			}, {
				Name:  "value1",
				Value: values[1],
			}, {
				Name:  "value2",
				Value: values[2],
			}, {
				Name:  "value3",
				Value: values[3],
			},
		}
		rowValues[0] = taosUtil.RowValue{Fields: fieldValues}
		subTableValue.Values = rowValues
		result = append(result, subTableValue)
	}
	return result
}
func val2float(data []byte) []string {
	length := len(data) / 4
	strings := make([]string, length)
	for i := 0; i < length; i++ {
		bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
		strings[i] = strconv.FormatFloat(float64(math.Float32frombits(bits)), 'f', 5, 64)
	}
	return splitArray(strings, 4)
}

func splitArray(source []string, n int64) []string {
	sLen := int64(len(source))
	var result = make([]string, n)
	remainder := sLen % n
	number := sLen / n
	var offset int64 = 0 //偏移量
	var i int64 = 0
	for ; i < n; i++ {
		var value []string
		if remainder > 0 {
			value = source[i*number+offset : (i+1)*number+offset+1]
			remainder--
			offset++
		} else {
			value = source[i*number+offset : (i+1)*number+offset]
		}
		join := strings.Join(value, ",")
		result[i] = join
	}
	return result
}

func test(data []byte) []float32 {
	length := len(data) / 4
	var result = make([]float32, length)
	for i := 0; i < length; i++ {
		bits := binary.LittleEndian.Uint32(data[i*4 : (i+1)*4])
		result[i] = math.Float32frombits(bits)
	}
	return result
}

func valSplitChannel(value string, count int) map[int][]byte {
	result := make(map[int][]byte)
	data := []byte(value)
	dataSum := len(data) / 4
	size := dataSum / count
	for i := 0; i < count; i++ {
		array := make([]byte, 0)
		result[i] = array
	}
	for j := 0; j < size; j++ {
		item := data[j*(count*4) : (j+1)*(count*4)]
		for a := 0; a < count; a++ {
			resultArray := result[a]
			resultArray = append(resultArray, item[a*4:(a+1)*4]...)
			result[a] = resultArray
		}
	}
	return result
}
func splitChannel(data []float32, count int) map[int][]float32 {
	result := make(map[int][]float32)
	size := len(data) / count
	for i := 0; i < count; i++ {
		array := make([]float32, 0)
		result[i] = array
	}
	for j := 0; j < size; j++ {
		item := data[j*count : (j+1)*count]
		for a := 0; a < count; a++ {
			resultArray := result[a]
			resultArray = append(resultArray, item[a])
			result[a] = resultArray
		}
	}
	return result
}

func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
