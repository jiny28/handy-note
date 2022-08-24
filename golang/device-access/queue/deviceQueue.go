package queue

import (
	"device-access/entity"
	"device-access/redisUtil"
	"device-access/taosUtil"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

var (
	taosInfo = taosUtil.TaosInfo{
		HostName:   "taos-server",
		ServerPort: 6030,
		User:       "root",
		Password:   "taosdata",
		DbName:     "hlhz1",
	}
	EventQueue     = make(chan entity.DeviceReceiveBean, 10000)
	batchSize      = 250
	workers        = 4
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
		funStart := time.Now()
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
				type SelfData struct {
					Device   string
					ItemCode string `json:"item_code"`
					Value    string
				}
				type SelfJson struct {
					Time entity.Time
					Data []SelfData
				}
				for _, bean := range v {
					payload := bean.Payload
					if payload == "" {
						continue
					}
					var jsonObject SelfJson
					err := json.Unmarshal([]byte(payload), &jsonObject)
					if err != nil {
						fmt.Printf("解析json错误,device:%v,topic:%v,error:%v", bean.Device, bean.Topic, err.Error())
						continue
					}
					time := time.Time(jsonObject.Time)
					payloadArray := jsonObject.Data
					for _, data := range payloadArray {
						device := data.Device
						code := data.ItemCode
						value := data.Value
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
				taosData := make(map[string]map[int64]map[string]interface{})
				fmt.Println("数据本身时间:", results[0].Time.Format("2006-01-02 15:04:05"), "当前时间为：", funStart.Format("2006-01-02 15:04:05"))
				for _, bean := range results {
					device := bean.Device
					code := bean.ItemCode
					time := bean.Time
					value := bean.Value
					redisData[device+"_"+code] = getRedisData(time, value)
					timeLong := time.UnixNano() / 1e6
					v, ok := taosData[device]
					tmap := make(map[int64]map[string]interface{})
					if ok {
						tmap = v
						vv, ok := tmap[timeLong]
						var vmap = make(map[string]interface{})
						if ok {
							vmap = vv
						}
						vmap[code] = value
						tmap[timeLong] = vmap
					} else {
						var vmap = make(map[string]interface{})
						vmap[code] = value
						tmap[timeLong] = vmap
					}
					taosData[device] = tmap
				}
				subTableValue := gerSubTableValue(taosData)
				fmt.Printf("解析过程耗时(batch:%v) = %v\n", batchSize, time.Since(funStart))
				taosStart := time.Now()
				_, e := taosUtil.InsertAutoCreateTable(subTableValue)
				fmt.Printf("taos insert 耗时(batch:%v) = %v\n", batchSize, time.Since(taosStart))
				if e != nil {
					fmt.Printf("taos insert error:" + e.Error())
				}
				redisStart := time.Now()
				rc.BatchSet(0, redisData, 0)
				fmt.Printf("redis insert 耗时(batch:%v) = %v\n", batchSize, time.Since(redisStart))
			} else if k == 2 {

			}
		}
		fmt.Printf("整个函数耗时 = %v\n", time.Since(funStart))
		return
	}
	errHandler = func(err error, batch []entity.DeviceReceiveBean) {
		fmt.Println("device queue error : ", err.Error())
	}
	getRedisData = func(time time.Time, value interface{}) string {
		type JsonObject struct {
			Date  string `json:"v_date"`
			Value string `json:"v_value"`
		}
		redisValue := JsonObject{Date: time.Format("2006-01-02 15:04:05"), Value: value.(string)}
		marshal, _ := json.Marshal(redisValue)
		return string(marshal)
	}
	gerSubTableValue = func(data map[string]map[int64]map[string]interface{}) []taosUtil.SubTableValue {
		var result = make([]taosUtil.SubTableValue, 0)
		for device, v := range data {
			var subTableValue taosUtil.SubTableValue
			subTableValue.Name = device + "_sub"
			subTableValue.SuperTable = device
			tags := []taosUtil.TagValue{
				{
					Name:  "device",
					Value: device,
				},
			}
			subTableValue.Tags = tags
			rowValues := make([]taosUtil.RowValue, len(v))
			var num = 0
			for time, tv := range v {
				fieldValues := make([]taosUtil.FieldValue, 0)
				fieldValues = append(fieldValues, taosUtil.FieldValue{
					Name:  "ts",
					Value: time,
				})
				for ik, iv := range tv {
					fieldValues = append(fieldValues, taosUtil.FieldValue{
						Name:  ik,
						Value: iv,
					})
				}
				rowValues[num] = taosUtil.RowValue{Fields: fieldValues}
				num++
			}
			subTableValue.Values = rowValues
			result = append(result, subTableValue)
		}
		return result
	}
)

func init() {
	taosUtil.Connection(taosInfo)
}

func RunDeviceQueue() {
	for i := 0; i < workers; i++ {
		go func() {
			/*for {
				var batch []entity.DeviceReceiveBean
				lens := len(EventQueue)
				if lens > batchSize {
					lens = batchSize
				}
				for o := 0; o < lens; o++ {
					msg := <-EventQueue
					batch = append(batch, msg)
				}
				if err := batchProcessor(batch); err != nil {
					errHandler(err, batch)
				}
				time.Sleep(time.Millisecond * 200)
			}*/
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
