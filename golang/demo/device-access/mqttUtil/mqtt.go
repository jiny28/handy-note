package mqttUtil

import (
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

//操作堵塞时间. s
var actionTimeout int = 1

type MqttConnection struct {
	Host               []string
	Client             string
	Username           string
	Password           string
	AutomaticReconnect bool
	CleanSession       bool
	connClient         mqtt.Client
}

func (m *MqttConnection) Connection(f mqtt.MessageHandler) {
	//mqtt.DEBUG = log.New(os.Stdout, "", 0)
	//mqtt.ERROR = log.New(os.Stdout, "", 0)
	opts := mqtt.NewClientOptions()
	for _, v := range m.Host {
		opts.AddBroker(v)
	}
	opts.SetClientID(m.Client)
	opts.SetUsername(m.Username)
	opts.SetPassword(m.Password)
	opts.SetAutoReconnect(m.AutomaticReconnect)
	opts.SetCleanSession(m.CleanSession)
	opts.SetKeepAlive(60 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	// 设置消息回调处理函数
	opts.SetDefaultPublishHandler(f)
	m.connClient = mqtt.NewClient(opts)
	if token := m.connClient.Connect(); token.WaitTimeout(time.Duration(actionTimeout)*time.Second) && token.Error() != nil {
		panic(token.Error())

	}
	mqtt.DEBUG.Println("mqtt connection success")
}

func (m *MqttConnection) Subscribe(topics map[string]byte, callback mqtt.MessageHandler) error {
	token := m.connClient.SubscribeMultiple(topics, callback)
	if token.WaitTimeout(time.Duration(actionTimeout)*time.Second) && token.Error() != nil {
		return token.Error()
	}
	mqtt.DEBUG.Println("subscribe topics:", topics)
	return nil
}

func (m *MqttConnection) UnSubscribe(topics ...string) error {
	if token := m.connClient.Unsubscribe(topics...); token.WaitTimeout(time.Duration(actionTimeout)*time.Second) && token.Error() != nil {
		return token.Error()
	}
	mqtt.DEBUG.Println("unsubscribe topic :", topics)
	return nil
}

func (m *MqttConnection) PublishMsg(topic string, qos byte, retain bool, payload interface{}) error {
	if token := m.connClient.Publish(topic, qos, retain, payload); token.Error() != nil {
		return token.Error()
	}
	mqtt.DEBUG.Println("publish msg on topic:", topic)
	return nil
}

func (m *MqttConnection) Disconnection(u uint) {
	m.connClient.Disconnect(u)
	mqtt.DEBUG.Println("disconnection.")
}
