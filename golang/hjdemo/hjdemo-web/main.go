package main

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"hjdemo-web/mqttUtil"
	"strconv"
	"strings"
	"time"
)

var mqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://192.168.1.229:1883"},
	Client:             "hjdemo-web",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       false,
}
var topics = map[string]byte{"exdevice/#": 0}

func init() {

}

func main() {
	mqttConnection.Connection(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		if payload == "" {
			return
		}
		topic := msg.Topic()
		//device := strings.Split(topic, "/")[1]
		ts := strings.Split(topic, "/")[2]
		//ts2 := strings.Split(payload, ",")[1]
		parseInt, _ := strconv.ParseInt(ts, 10, 64)
		//parseInt2, _ := strconv.ParseInt(ts2, 10, 64)
		nowMs := time.Now().UnixNano() / 1e6
		fmt.Printf("上行链路时间为(ms)：%v\n", nowMs-parseInt)
		//fmt.Printf("完整链路时间为(ms)：%v\n", nowMs-parseInt2)
		//fmt.Printf("接收到%v设备的数据,数据的大小为:%v\n", device, len(payload))

	})
	err := mqttConnection.Subscribe(topics, nil)
	checkErr(err, "mqtt sub err:")
	for {
		time.Sleep(time.Second * 60 * 10)
	}
	mqttConnection.Disconnection(250)
}
func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
