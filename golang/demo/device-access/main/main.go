package main

import (
	"database/sql"
	"device-access/mqttUtil"
	"device-access/mysqlUtil"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var mysqlInfo = mysqlUtil.MysqlInfo{
	UserName: "root",
	Password: "123456",
	Ip:       "127.0.0.1",
	Port:     33066,
	Db:       "test",
}

var mqttConnection = mqttUtil.MqttConnection{
	Host:               []string{"tcp://10.88.0.14:1883"},
	Client:             "go_admin",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       true,
}

func main() {
	db := initMysql()
	defer db.Close()
	exeSql := "select * from test"
	res, err := mysqlUtil.GetAll(exeSql)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("result : ", res)
	mqttConnection.Connection(f)
	err = mqttConnection.Subscribe(map[string]byte{"aaa": 0, "bbb": 0}, nil)
	if err != nil {
		panic(err.Error())
	}
	err = mqttConnection.PublishMsg("test", 0, false, "发送消息")
	if err != nil {
		panic(err.Error())
	}
	for true {
		time.Sleep(1000)
	}
}

func initMysql() *sql.DB {
	db, error := mysqlUtil.Connection(mysqlInfo)
	if error != nil {
		panic(error.Error())
	}
	fmt.Println(" mysql connection success")
	return db
}
