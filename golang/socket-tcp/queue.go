package main

import (
	"encoding/json"
	"fmt"
	"socket-tcp/taosUtil"
	"strconv"
	"strings"
	"time"
)

var EventQueue = make(chan map[string][]string, 10000)
var workers = 2
var batchSize = 5

func RunQueue() {
	go func() {
		for {
			lens := len(EventQueue)
			if lens == 0 {
				time.Sleep(200 * time.Millisecond)
				continue
			}
			if lens > batchSize {
				lens = batchSize
			}
			batch := make([]map[string][]string, 0)
			for o := 0; o < lens; o++ {
				msg := <-EventQueue
				batch = append(batch, msg)
			}
			go batchProcessor(batch)
		}
	}()
}
func batchProcessor(batch []map[string][]string) {
	fmt.Printf("接收到数据：%v \n", len(batch))
	redisRes := make(map[string][]float64)
	for _, v := range batch {
		for device, oneGroupData := range v {
			itemInfoArr := itemMapping[device]
			size := len(oneGroupData)
			var ts int64
			for i := 1; i <= size; i++ {
				data := oneGroupData[i-1]
				split := strings.Split(data, ",")
				if ts == 0 {
					ts, _ = strconv.ParseInt(split[0], 10, 64)
				}
				for j := 1; j <= len(split)-1; j++ { //遍历一条字符数据 获取采集值
					itemInfoMap := itemInfoArr[j-1]
					itemCode := itemInfoMap["item_code"]
					value := split[j]
					dvalue, _ := strconv.ParseFloat(value, 64)
					reList, ok := redisRes[device+"_"+itemCode+"_"+fmt.Sprintf("%v", ts)]
					if !ok {
						reList = make([]float64, 0)

					}
					reList = append(reList, dvalue)
					redisRes[device+"_"+itemCode+"_"+fmt.Sprintf("%v", ts)] = reList
				}
			}
		}
	}
	taos := convertTaos(redisRes)
	_, err := taosUtil.InsertAutoCreateTable(taos)
	checkErr(err, "insert taos error")
}

func convertTaos(redisRes map[string][]float64) []taosUtil.SubTableValue {
	result := make([]taosUtil.SubTableValue, 0)
	for k, objects := range redisRes {
		splitK := strings.Split(k, "_")
		deviceCode := splitK[0]
		itemCode := splitK[1]
		time := splitK[2]
		ts, _ := strconv.ParseInt(time, 10, 64)
		tags := []taosUtil.TagValue{
			{
				Name:  "device",
				Value: deviceCode,
			},
			{
				Name:  "item_code",
				Value: itemCode,
			},
		}
		marshal, _ := json.Marshal(objects)
		values := []taosUtil.RowValue{{
			Fields: []taosUtil.FieldValue{
				{
					Name:  "ts",
					Value: ts,
				}, {
					Name:  "value",
					Value: string(marshal),
				},
			},
		}}
		var subTableValue taosUtil.SubTableValue
		subTableValue.Name = deviceCode + "_" + itemCode
		subTableValue.SuperTable = "meter1"
		subTableValue.Tags = tags
		subTableValue.Values = values
		result = append(result, subTableValue)
	}
	return result
}
