package main

import (
	"database/sql"
	"device-access/entity"
	"device-access/mqttUtil"
	"device-access/mysqlUtil"
	"device-access/queue"
	"device-access/redisUtil"
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"strconv"
	"time"
)

var mysqlInfo = mysqlUtil.MysqlInfo{
	UserName: "root",
	Password: "123456",
	Ip:       "10.88.0.14",
	Port:     33066,
	Db:       "hlhz_go_test",
}

var mqttConnection = mqttUtil.MqttConnection{
	//Host:               []string{"tcp://10.88.0.14:1883"},
	Client:             "go_admin",
	Username:           "hlhz",
	Password:           "hlhz.123456",
	AutomaticReconnect: true,
	CleanSession:       true,
}

var redisInfo = redisUtil.RedisInfo{
	Ip:       "10.88.0.14",
	Port:     63799,
	Password: "xapp",
	Db:       0,
}

var db *sql.DB

func init() {
	redisUtil.RedisInit(redisInfo)
	db = initMysql()
}

func main() {
	fmt.Println("设备接入启动.")
	defer db.Close()
	mqttInfos, err := mysqlUtil.GetAll("SELECT c_ip,c_port FROM t_mqtt_info WHERE c_state = 1")
	checkErr(err, "get mqtt info sql")
	if len(mqttInfos) == 0 {
		fmt.Println("无启用的mqtt服务器.")
		return
	}
	mqttTopics, err := mysqlUtil.GetAll("SELECT c_topic,c_device,c_method_index FROM t_mqtt_topic WHERE c_state = 1")
	if len(mqttTopics) == 0 {
		fmt.Println("无启用的主题.")
		return
	}
	mqttAddrs := convertMqttAddr(mqttInfos)
	topics := make(map[string]byte)
	topicByIndex := make(map[string]int)
	topicByDevice := make(map[string]string)
	for _, topicMap := range mqttTopics {
		top := topicMap["c_topic"]
		topics[top] = 2
		methodIndex, _ := strconv.Atoi(topicMap["c_method_index"])
		topicByIndex[top] = methodIndex
		topicByDevice[top] = topicMap["c_device"]
	}
	queue.RunDeviceQueue()
	mqttConnection.Host = mqttAddrs
	mqttConnection.Connection(func(client mqtt.Client, msg mqtt.Message) {
		payload := string(msg.Payload())
		if payload == "" {
			return
		}
		topic := msg.Topic()
		index, ok := topicByIndex[topic]
		if !ok {
			return
		}
		device := topicByDevice[topic]
		bean := entity.DeviceReceiveBean{Topic: topic, Device: device, MethodIndex: index, Payload: payload}
		queue.EventQueue <- bean
	})
	defer mqttConnection.Disconnection(250)
	err = mqttConnection.Subscribe(topics, nil)
	rc := redisUtil.Redis{}
	for {
		time.Sleep(time.Second * 60)
		flag, e := rc.Get(10, "deviceReceiveConfigFlag")
		checkErr(e, "redis get deviceReceiveConfigFlag")
		if flag == "1" {
			er := rc.Set(10, "deviceReceiveConfigFlag", "0", 0)
			checkErr(er, "redis set deviceReceiveConfigFlag")
			break
		}
	}
	//test()
}

func test() {
	/*rc := redisUtil.Redis{}
	flag, e := rc.Get(10, "deviceReceiveConfigFlag")
	checkErr(e, "redis get deviceReceiveConfigFlag")
	if flag == "1" {
		fmt.Println("true")
	}*/
}

func convertMqttAddr(infos []map[string]string) []string {
	result := make([]string, len(infos))
	for i, info := range infos {
		result[i] = "tcp://" + info["c_ip"] + ":" + info["c_port"]
	}
	return result
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
