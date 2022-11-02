package main

import (
	"encoding/json"
	"fmt"
	"jkRedis/redisUtil"
	"jkRedis/taosUtil"
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
var redisInfo = redisUtil.RedisInfo{
	Ip:       "127.0.0.1",
	Port:     63799,
	Password: "xapp",
	Db:       0,
}
var rc redisUtil.Redis

func main() {
	for {
		//startTime := time.Now()
		items := []string{"yali", "shidu", "wendu"}
		result, err := rc.MGet(13, items)
		if err != nil {
			fmt.Println("redis error : " + err.Error())
			return
		}
		type SelfData struct {
			V float32
			Q int
			T int64
		}
		timeMap := make(map[int64][]string, 0)
		for i, v := range result {
			if v == nil {
				fmt.Printf("下标为:%v没有找到值.\n", i)
				continue
			}
			var selfData SelfData
			err := json.Unmarshal([]byte(v.(string)), &selfData)
			if err != nil {
				fmt.Printf("解析json错误,key下标 : %v , error :%v \n", i, err.Error())
				continue
			}
			//fmt.Printf("key下标:%v,t:%v,v:%v\n", i, selfData.T, selfData.V)
			array, ok := timeMap[selfData.T]
			fieldAndValue := items[i] + "," + fmt.Sprintf("%.2f", selfData.V)
			if !ok {
				array = []string{fieldAndValue}
			} else {
				array = append(array, fieldAndValue)
			}
			timeMap[selfData.T] = array
		}
		if len(timeMap) == 0 {
			fmt.Println("redis not value .")
			continue
		}
		writeTaos(timeMap)
		//fmt.Printf("全程耗时:%v\n", time.Since(startTime))
		time.Sleep(time.Millisecond * 1000)
	}
}

func writeTaos(datas map[int64][]string) {
	rowValues := make([]taosUtil.RowValue, 0)
	for time, values := range datas {
		time = time * 1000 // to us
		fieldValues := []taosUtil.FieldValue{
			{
				Name:  "ts",
				Value: time,
			},
		}
		for _, v := range values {
			split := strings.Split(v, ",")
			name := split[0]
			value := split[1]
			fieldValues = append(fieldValues, taosUtil.FieldValue{Name: name, Value: value})
		}
		rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
	}
	var subTableValue taosUtil.SubTableValue
	subTableValue.Name = "jk0"
	subTableValue.SuperTable = "jk"
	tags := []taosUtil.TagValue{
		{
			Name:  "groupId",
			Value: "0",
		},
	}
	subTableValue.Tags = tags
	subTableValue.Values = rowValues
	_, err := taosUtil.InsertAutoCreateTable([]taosUtil.SubTableValue{subTableValue})
	if err != nil {
		fmt.Println("insert taos error :" + err.Error())
	}
}
func init() {
	redisUtil.RedisInit(redisInfo)
	taosUtil.Connection(taosInfo)
	rc = redisUtil.Redis{}
}
