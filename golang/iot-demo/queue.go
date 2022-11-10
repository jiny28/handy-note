package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"iot-demo/entity"
	"iot-demo/taosUtil"
	"time"
)

var EventQueue = make(chan entity.DeviceReceiveBean, 1000000)
var workers = 2
var batchSize = 30

func runDeviceQueue() {
	//for i := 0; i < workers; i++ {
	go func() {
		for {
			lens := len(EventQueue)
			if lens == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			if lens > batchSize {
				lens = batchSize
			}
			batch := make([]entity.DeviceReceiveBean, 0)
			for o := 0; o < lens; o++ {
				msg := <-EventQueue
				batch = append(batch, msg)
			}
			go batchProcessor(batch)
		}
	}()
	//}
}
func batchProcessor(batch []entity.DeviceReceiveBean) {
	fmt.Printf("接收到数据：%v \n", len(batch))
	result := make([]taosUtil.SubTableValue, 0)
	for _, obj := range batch {
		tags := []taosUtil.TagValue{
			{
				Name:  "device",
				Value: obj.Device,
			},
		}
		payload := obj.Payload
		if payload == "" {
			continue
		}
		hexData, _ := hex.DecodeString(payload)
		payload = string(hexData)
		var jsonObject entity.SelfJson
		err := json.Unmarshal([]byte(payload), &jsonObject)
		if err != nil {
			fmt.Printf("解析json错误,device:%v,topic:%v,error:%v", obj.Device, obj.Topic, err.Error())
			continue
		}
		ts := jsonObject.Time
		data := jsonObject.Data
		rowValues := make([]taosUtil.RowValue, 0)
		fieldValues := make([]taosUtil.FieldValue, 0)
		fieldValues = append(fieldValues, taosUtil.FieldValue{
			Name:  "ts",
			Value: ts,
		})
		for _, m := range data {
			for k, v := range m {
				fieldValues = append(fieldValues, taosUtil.FieldValue{
					Name:  k,
					Value: v,
				})
			}
		}
		rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
		var subTableValue taosUtil.SubTableValue
		subTableValue.Name = obj.Device
		subTableValue.SuperTable = "kyj"
		subTableValue.Tags = tags
		subTableValue.Values = rowValues
		result = append(result, subTableValue)
	}
	_, err := taosUtil.InsertAutoCreateTable(result)
	if err != nil {
		fmt.Println("taos insert error :" + err.Error())
		panic(err.Error())
	}
}
