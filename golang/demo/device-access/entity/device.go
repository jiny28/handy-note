package entity

import "time"

type DeviceReceiveBean struct {
	// 接收到的主题
	Topic string
	// 主题对应的设备编码，是否为空根据业务来，有的主题的消息内容没有设备编码，那就是配置的这里的设备编码
	Device string
	// 解析方法
	MethodIndex int
	// 数据
	Payload string
}

type DeviceStandardBean struct {
	Time     time.Time
	Device   string
	ItemCode string
	Value    interface{}
}
