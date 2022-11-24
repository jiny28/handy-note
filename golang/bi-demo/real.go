package main

import (
	"bi-demo/entity"
	"bi-demo/mqttUtil"
	"encoding/json"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strings"
	"sync"
	"time"
)

var mqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://127.0.0.1:1883"},
	Client:             "bi",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       false,
}
var topics = map[string]byte{"exiot/#": 2, "deal/#": 2}
var EventQueue = make(chan entity.DeviceReceiveBean, 1000000)
var RealValue sync.Map //<key(device-itemCode)-{value,time}>
var layout = "2006-01-02 15:04:05"

func startSub() {
	mqttConnection.Connection(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		if payload == "" {
			return
		}
		topic := msg.Topic()
		bean := entity.DeviceReceiveBean{Topic: topic, Payload: payload}
		EventQueue <- bean
	})
	err := mqttConnection.Subscribe(topics, nil)
	checkErr(err, "mqtt sub err:")
	runDeviceQueue()
}

func runDeviceQueue() {
	go func() {
		for {
			select {
			case msg := <-EventQueue:
				topic := msg.Topic
				if strings.Contains(topic, "exiot/") {
					device := strings.Split(topic, "/")[1]
					payload := msg.Payload
					var jsonObject entity.SelfJson
					err := json.Unmarshal([]byte(payload), &jsonObject)
					if err != nil {
						fmt.Printf("解析json错误,device:%v,error:%v", device, err.Error())
						continue
					}
					ts := jsonObject.Time
					data := jsonObject.Data
					strTime := time.Unix(ts/1000, 0).Format(layout)
					for _, v := range data {
						for vk, vv := range v {
							key := device + "-" + vk
							RealValue.Store(key, entity.RealValue{Date: strTime, Value: vv})
						}
					}
				} else if strings.Contains(topic, "deal/") {
					key := strings.Split(topic, "/")[1]
					payload := msg.Payload
					type PayloadStruct struct {
						Time  string
						Value string
					}
					var jsonObject PayloadStruct
					err := json.Unmarshal([]byte(payload), &jsonObject)
					if err != nil {
						fmt.Printf("解析json错误,topic:%v,error:%v", topic, err.Error())
						continue
					}
					RealValue.Store(key, entity.RealValue{Date: jsonObject.Time, Value: jsonObject.Value})
				}
			}
		}
	}()
}
