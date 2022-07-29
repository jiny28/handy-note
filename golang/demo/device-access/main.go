package main

import (
	"database/sql"
	"device-access/mqttUtil"
	"device-access/mysqlUtil"
	"device-access/redisUtil"
	"device-access/taosUtil"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"math/rand"
	"strconv"
	"time"
)

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

var f mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	fmt.Printf("TOPIC: %s\n", msg.Topic())
	fmt.Printf("MSG: %s\n", msg.Payload())
}

var redisInfo = redisUtil.RedisInfo{
	Ip:       "10.88.0.14",
	Port:     63799,
	Password: "xapp",
	Db:       0,
}

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "hlhz",
}

func main() {
	db := initMysql()
	defer db.Close()
	exeSql := "select * from test"
	res, err := mysqlUtil.GetAll(exeSql)
	checkErr(err, "sql:"+exeSql)
	fmt.Println("result : ", res)
	mqttConnection.Connection(f)
	defer mqttConnection.Disconnection(250)
	err = mqttConnection.Subscribe(map[string]byte{"aaa": 0, "bbb": 0}, nil)
	checkErr(err, "subscribe topic ")
	err = mqttConnection.PublishMsg("test", 0, false, "发送消息")
	checkErr(err, "PublishMsg")
	/*for true {
		time.Sleep(1000)
	}*/
	redisUtil.RedisInit(redisInfo)
	redisUtil.RedisClient.Do("select", 2)
	err = redisUtil.RedisClient.Set("aaa", 123, 0).Err()
	checkErr(err, "redis set ")
	result, errGet := redisUtil.RedisClient.Get("aaa").Result()
	checkErr(errGet, "redis get ")
	fmt.Println(result)
	_, err = taosUtil.Connection(taosInfo)
	checkErr(err, "taos connection")
	defer taosUtil.Close()
	data := createTaosData()
	num, err := taosUtil.InsertAutoCreateTable(data)
	checkErr(err, "taos insert")
	fmt.Println("插入多少条", num)

}

func createTaosData() []taosUtil.SubTableValue {
	var subTableValues []taosUtil.SubTableValue
	startTime := "2022-07-10 10:00:00"

	startDate, err := time.ParseInLocation("2006-01-02 15:04:05", startTime, time.Local)

	checkErr(err, "")
	start := startDate.UnixNano()
	// 100个电表，每块电表200个点位
	for i := 0; i < 1; i++ {
		var subTableValue taosUtil.SubTableValue
		is := strconv.Itoa(i)
		subTableValue.Name = "d00" + is
		subTableValue.SuperTable = "meters"
		tags := []taosUtil.TagValue{
			{
				Name:  "location",
				Value: "d00" + is,
			},
			{
				Name:  "groupId",
				Value: 1 + i,
			},
		}
		subTableValue.Tags = tags
		num := 100 // 一次插入多少数据
		rowValues := make([]taosUtil.RowValue, num)
		for j := 0; j < num; j++ {
			fieldNum := 200
			fieldValues := make([]taosUtil.FieldValue, fieldNum)
			fieldValues[0] = taosUtil.FieldValue{
				Name:  "ts",
				Value: start,
			}
			for a := 0; a < fieldNum; a++ {

				fieldName := "field" + strconv.Itoa(a)

				f := rand.Int() % 1000
				v := 200 + rand.Float32()

				var fieldValue taosUtil.FieldValue
				if a%2 == 0 {
					fieldValue = taosUtil.FieldValue{
						Name:  fieldName,
						Value: v,
					}
				} else {
					fieldValue = taosUtil.FieldValue{
						Name:  fieldName,
						Value: f,
					}
				}
				fieldValues = append(fieldValues, fieldValue)
			}
			rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
			start += 1000
		}
		subTableValue.Values = rowValues
		subTableValues = append(subTableValues, subTableValue)
	}
	return subTableValues

}

func initMysql() *sql.DB {
	db, error := mysqlUtil.Connection(mysqlInfo)
	checkErr(error, "初始化MySql错误")
	fmt.Println(" mysql connection success")
	return db
}

func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
