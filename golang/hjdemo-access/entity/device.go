package entity

import "time"

type DeviceReceiveBean struct {
	// 接收到的主题
	Topic string
	// 主题对应的设备编码，是否为空根据业务来，有的主题的消息内容没有设备编码，那就是配置的这里的设备编码
	Device string
	// 解析方法
	//MethodIndex int
	// 数据
	Payload string
}

type SelfData struct {
	ItemCode string      `json:"itemCode"`
	Value    interface{} `json:"value"`
}
type SelfJson struct {
	Time int64                    `json:"time"`
	Data []map[string]interface{} `json:"data"`
}

type DeviceStandardBean struct {
	Time     time.Time
	Device   string
	ItemCode string
	Value    interface{}
}

const (
	timeFormat = "2006-01-02 15:04:05"
)

type Time time.Time

func (t *Time) UnmarshalJSON(data []byte) (err error) {
	now, err := time.ParseInLocation(`"`+timeFormat+`"`, string(data), time.Local)
	*t = Time(now)
	return
}

func (t Time) MarshalJSON() ([]byte, error) {
	b := make([]byte, 0, len(timeFormat)+2)
	b = append(b, '"')
	b = time.Time(t).AppendFormat(b, timeFormat)
	b = append(b, '"')
	return b, nil
}

func (t Time) String() string {
	return time.Time(t).Format(timeFormat)
}
