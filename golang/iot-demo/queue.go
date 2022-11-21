package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"iot-demo/entity"
	"iot-demo/gopool"
	"iot-demo/mqttUtil"
	"iot-demo/taosUtil"
	"time"
)

var EventQueue = make(chan entity.DeviceReceiveBean, 1000000)
var batchSize = 30
var inter = 100 * time.Millisecond
var poolNum = 2
var jobQueueNum = 5
var workerPool *gopool.WorkerPool

var externalMqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://127.0.0.1:1883"},
	Client:             "iot-exter",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       false,
}

func init() {
	fmt.Printf("协程池初始化:poolNum:%v,jobQueueNum:%v\n", poolNum, jobQueueNum)
	workerPool = gopool.NewWorkerPool(poolNum, jobQueueNum)
	workerPool.Start()
	externalMqttConnection.Connection(func(client mqtt.Client, message mqtt.Message) {
		fmt.Println("external mqtt print msg:" + string(message.Payload()))
	})
}

type Task struct {
	batch []entity.DeviceReceiveBean
}

func runDeviceQueue() {
	go func() {
		for {
			lens := len(EventQueue)
			fmt.Printf("当前队列大小:%v\n", lens)
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
			tJob := Task{batch: batch}
			workerPool.JobQueue <- tJob
			time.Sleep(inter)
		}
		externalMqttConnection.Disconnection(250)
	}()
}

func (t Task) RunTask(request interface{}) {
	batchProcessor(t.batch)
}

func batchProcessor(batch []entity.DeviceReceiveBean) {
	//fmt.Printf("接收到数据：%v \n", len(batch))
	//startNow := time.Now()
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
		mqttError := externalMqttConnection.PublishMsg("exiot/"+obj.Device, 0, false, payload)
		if mqttError != nil {
			fmt.Printf("mqtt转发错误device:%v:%v\n", obj.Device, mqttError.Error())
			continue
		}
		var jsonObject entity.SelfJson
		err := json.Unmarshal([]byte(payload), &jsonObject)
		if err != nil {
			fmt.Printf("解析json错误,device:%v,topic:%v,error:%v", obj.Device, obj.Topic, err.Error())
			continue
		}
		ts := jsonObject.Time
		//fmt.Println("处理数据的时间为：" + time.Unix(0, ts*int64(time.Millisecond)).Format("2006-02-01 15:04:05.000"))
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
	//startTaos := time.Now()
	_, err := taosUtil.InsertAutoCreateTable(result)
	if err != nil {
		fmt.Println("taos insert error :" + err.Error())
		panic(err.Error())
	}
	//fmt.Printf("save taos 耗时:%v\n", time.Since(startTaos))
	//fmt.Printf("batchProcessor 耗时:%v\n", time.Since(startNow))
}
