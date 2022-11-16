package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"hjdemo-access/mqttUtil"
	"hjdemo-access/taosUtil"
	"strconv"
	"strings"
	"time"
)

var EventQueue = make(chan string, 1000000)
var batchSize = 25
var inter = 10 * time.Millisecond
var poolNum = 5
var jobQueueNum = 5
var workerPool *WorkerPool

var mqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://127.0.0.1:1883"},
	Client:             "hjdemo_access",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       false,
}

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "h",
}

func init() {
	fmt.Printf("协程池初始化:poolNum:%v,jobQueueNum:%v\n", poolNum, jobQueueNum)
	workerPool = NewWorkerPool(poolNum, jobQueueNum)
	workerPool.Start()
	mqttConnection.Connection(func(client mqtt.Client, message mqtt.Message) {
		fmt.Println("external mqtt print msg:" + string(message.Payload()))
	})
	taosUtil.Connection(taosInfo)
}

type Task struct {
	batch []string
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
			batch := make([]string, 0)
			for o := 0; o < lens; o++ {
				msg := <-EventQueue
				batch = append(batch, msg)
			}
			tJob := Task{batch: batch}
			workerPool.JobQueue <- tJob
			time.Sleep(inter)
		}
		mqttConnection.Disconnection(250)
	}()
}

func (t Task) RunTask(request interface{}) {
	batchProcessor(t.batch)
}

func batchProcessor(batch []string) {
	fmt.Printf("接收到数据大小：%v \n", len(batch[0]))
	startNow := time.Now()
	result := make([]taosUtil.SubTableValue, 0)
	for _, obj := range batch {
		split := strings.Split(obj, ",")
		device := split[0]
		tags := []taosUtil.TagValue{
			{
				Name:  "device",
				Value: device,
			},
		}
		mqttError := mqttConnection.PublishMsg("exdevice/"+device, 0, false, obj)
		if mqttError != nil {
			fmt.Printf("mqtt转发错误device:%v:%v\n", device, mqttError.Error())
			continue
		}
		ts := time.Now().UnixNano() / 1e6
		rowValues := make([]taosUtil.RowValue, 0)
		fieldValues := make([]taosUtil.FieldValue, 0)
		fieldValues = append(fieldValues, taosUtil.FieldValue{
			Name:  "ts",
			Value: ts,
		})
		for i := 1; i < len(split); i++ {
			varName := "item" + strconv.Itoa(i-1)
			fieldValues = append(fieldValues, taosUtil.FieldValue{
				Name:  varName,
				Value: split[i],
			})
		}
		rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
		var subTableValue taosUtil.SubTableValue
		subTableValue.Name = device
		subTableValue.SuperTable = "hjdemo"
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
	fmt.Printf("batchProcessor 耗时:%v\n", time.Since(startNow))
}
