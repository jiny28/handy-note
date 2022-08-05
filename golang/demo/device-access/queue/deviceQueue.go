package queue

import (
	"device-access/entity"
	"device-access/redisUtil"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	EventQueue     = make(chan entity.DeviceReceiveBean, 10000)
	batchSize      = 1000
	workers        = 1
	batchProcessor = func(batch []entity.DeviceReceiveBean) (e error) {
		defer func() {
			if err := recover(); err != nil {
				str, ok := err.(string)
				if ok {
					e = errors.New(str)
				} else {
					e = errors.New("panic")
				}
			}
		}()
		rc := redisUtil.Redis{}
		groupByMethod := make(map[int][]entity.DeviceReceiveBean)
		for _, bean := range batch {
			index := bean.MethodIndex
			beans, ok := groupByMethod[index]
			list := make([]entity.DeviceReceiveBean, 0)
			if ok {
				list = beans
			}
			list = append(list, bean)
			groupByMethod[index] = list
		}
		for k, v := range groupByMethod {
			if k == 1 {
				results := make([]entity.DeviceStandardBean, 0)
				type selfData struct {
					device    string
					item_code string
					value     string
				}
				type selfJson struct {
					time time.Time
					data []selfData
				}
				for _, bean := range v {
					payload := bean.Payload
					if payload == "" {
						continue
					}
					var jsonObject selfJson
					err := json.Unmarshal([]byte(payload), &jsonObject)
					if err != nil {
						fmt.Printf("解析json错误,device:%v,topic:%v,error:%v", bean.Device, bean.Topic, err.Error())
						continue
					}
					time := jsonObject.time
					payloadArray := jsonObject.data
					for _, data := range payloadArray {
						device := data.device
						code := data.item_code
						value := data.value
						if device == "" || code == "" || value == "" {
							continue
						}
						results = append(results, entity.DeviceStandardBean{
							Time: time, Device: device, ItemCode: code, Value: value,
						})
					}
				}
				if len(results) == 0 {
					return
				}
				redisData := make(map[string]interface{})
				for _, bean := range results {
					device := bean.Device
					code := bean.ItemCode
					time := bean.Time
					value := bean.Value
					redisData[device+"-"+code] = getRedisData(time, value)
				}
				rc.BatchSet(0, redisData, 0)
			} else if k == 2 {

			}
		}
		return
	}
	errHandler = func(err error, batch []entity.DeviceReceiveBean) {
		fmt.Println("device queue error : ", err.Error())
	}
	getRedisData = func(time time.Time, value interface{}) string {
		type jsonObject struct {
			v_data  string
			v_value string
		}
		marshal, _ := json.Marshal(jsonObject{time.Format("2006-01-02 15:04:05"), value.(string)})
		return string(marshal)
	}
)

func RunDeviceQueue() {
	for i := 0; i < workers; i++ {
		go func() {
			var batch []entity.DeviceReceiveBean
			for {
				select {
				case msg := <-EventQueue:
					batch = append(batch, msg)
					if len(batch) != batchSize {
						break
					}
					if err := batchProcessor(batch); err != nil {
						errHandler(err, batch)
					}
					batch = make([]entity.DeviceReceiveBean, 0)
				}
			}
		}()
	}
}
