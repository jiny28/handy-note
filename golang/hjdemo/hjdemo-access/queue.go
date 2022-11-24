package main

import (
	"bufio"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"hjdemo-access/mqttUtil"
	"hjdemo-access/taosUtil"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var EventQueue = make(chan string, 1000000)
var batchSize = 20
var inter = 10 * time.Millisecond
var poolNum = 5
var jobQueueNum = 5
var workerPool *WorkerPool

var mqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://192.168.1.229:1883"},
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
		var startTime time.Time
		flag := false
		for {
			lens := len(EventQueue)
			//fmt.Printf("当前队列大小:%v\n", lens)
			if lens > 0 && !flag {
				startTime = time.Now()
				flag = true
				fmt.Println("开始消费")
			}
			if lens == 0 && flag {
				// 一波数据结束
				fmt.Printf("一波时间数据消费完数据处理耗时:%v\n", dealSumTime)
				fmt.Printf("一波时间数据消费完emqx耗时:%v\n", mqttSumTime)
				fmt.Printf("一波事件数据消费完总耗时:%v\n", time.Since(startTime))
				flag = false
				dealSumTime = time.Duration(1) * time.Millisecond
				mqttSumTime = time.Duration(1) * time.Millisecond
			}
			if lens == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			/*if lens > batchSize {
				lens = batchSize
			}*/
			batch := make([]string, 0)
			for o := 0; o < lens; o++ {
				msg := <-EventQueue
				batch = append(batch, msg)
			}
			/*tJob := Task{batch: batch}
			workerPool.JobQueue <- tJob*/
			if lens > batchSize {
				//获取应该数组分割为多少份
				var quantity int64
				if lens%batchSize == 0 {
					quantity = int64(lens) / int64(batchSize)
				} else {
					quantity = (int64(lens) / int64(batchSize)) + 1
				}
				//声明分割数组的截止下标
				var start, end, i int64
				for i = 1; i <= quantity; i++ {
					end = i * int64(batchSize)
					if i != quantity {
						batchProcessor(batch[start:end])
					} else {
						batchProcessor(batch[start:])
					}
					start = i * int64(batchSize)
				}
			} else {
				batchProcessor(batch)
			}
			time.Sleep(inter)
		}
		mqttConnection.Disconnection(250)
	}()
}

func (t Task) RunTask(request interface{}) {
	batchProcessor(t.batch)
}

var wg sync.WaitGroup

func batchProcessor(batch []string) {
	//startNow := time.Now()
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
	//fmt.Printf("组装%v台设备的taos对象耗时:%v\n", len(batch), time.Since(startNow))
	startTaos := time.Now()
	_, err := taosUtil.InsertAutoCreateTable(result)
	if err != nil {
		fmt.Println("taos insert error :" + err.Error())
		//panic(err.Error())
	}
	fmt.Printf("存储%v台设备的taos对象耗时:%v\n", len(batch), time.Since(startTaos))
}

func batchProcessorThread(batch []string) {
	//startNow := time.Now()
	for a := 0; a < 3; a++ {
		wg.Add(1)
		go func(batch []string) {
			defer wg.Done()
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
			//fmt.Printf("组装%v台设备的taos对象耗时:%v\n", len(batch), time.Since(startNow))
			startTaos := time.Now()
			_, err := taosUtil.InsertAutoCreateTable(result)
			if err != nil {
				fmt.Println("taos insert error :" + err.Error())
				//panic(err.Error())
			}
			fmt.Printf("存储%v台设备的taos对象耗时:%v\n", len(batch), time.Since(startTaos))
		}(batch)
	}
	wg.Wait()
}

func writerFile(batch []string) {
	//创建一个新文件，写入内容 5 句 “http://c.biancheng.net/golang/”
	now := time.Now()
	filePath := "/var/golang.txt"
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		fmt.Println("文件打开失败", err)
	}
	//及时关闭file句柄
	defer file.Close()
	//写入文件时，使用带缓存的 *Writer
	write := bufio.NewWriter(file)
	for i := range batch {
		write.WriteString(batch[i] + "\n")
	}
	//Flush将缓存的文件真正写入到文件中
	write.Flush()
	fmt.Printf("写入%v台设备的数据至文件耗时:%v\n", len(batch), time.Since(now))
}
