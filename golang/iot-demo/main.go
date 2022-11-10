package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"iot-demo/entity"
	"iot-demo/mqttUtil"
	"iot-demo/taosUtil"
	"strings"
	"time"
)

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "h",
}
var mqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://10.88.0.14:1883"},
	Client:             "go_admin",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       false,
}

func init() {
	taosUtil.Connection(taosInfo)
}

var taskFlag = make(chan int, 10)
var RestartFlag = false
var topics = map[string]byte{"device/#": 2}

func main() {
	runDeviceQueue()
	go task()
	// start http 服务
	go startWeb()
	for {
		select {
		case <-taskFlag:
			fmt.Println("restart .")
			go task()
		}
	}
	fmt.Println("main over.")
}

func task() {
	fmt.Println("iot go run.")
	//初始化xml
	mqttConnection.Connection(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		if payload == "" {
			return
		}
		topic := msg.Topic()
		device := strings.Split(topic, "/")[1]
		bean := entity.DeviceReceiveBean{Topic: topic, Device: device, Payload: payload}
		EventQueue <- bean
	})
	err := mqttConnection.Subscribe(topics, nil)
	checkErr(err, "mqtt sub err:")
	for {
		time.Sleep(time.Second * 60)
		if RestartFlag {
			RestartFlag = false
			break
		}
	}
	mqttConnection.Disconnection(250)
	fmt.Println("iot go exit.")
	taskFlag <- 1
}

func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
