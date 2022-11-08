package main

import (
	"fmt"
	"iot-demo/taosUtil"
	"math/rand"
	"strconv"
	"time"
)

type MoniObj struct {
	airStatus          string
	airElectric        float64
	airGas             float64
	airWater           float64
	airEscapage        float64
	airFlow            int
	airPower           int
	airExhaustPressure float64
	airGasDisplacement float64
	airMainPressure    int
	ts                 int64
}

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "h",
}
var EventQueue = make(chan MoniObj, 100)
var batchSize = 50

func init() {
	taosUtil.Connection(taosInfo)
}

// CREATE STABLE device (ts timestamp, airStatus binary(64), airElectric float, airGas float, airWater float, airEscapage float, airFlow INT, airPower INT, airExhaustPressure float, airGasDisplacement float,airMainPressure INT) TAGS (device binary(64));

func main() {
	go startMoni()
	// start http 服务
	go startWeb()
	fmt.Println("数据开始接入.")
	for {
		lens := len(EventQueue)
		if lens == 0 {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if lens > batchSize {
			lens = batchSize
		}
		batch := make([]MoniObj, 0)
		for o := 0; o < lens; o++ {
			msg := <-EventQueue
			batch = append(batch, msg)
		}
		go batchProcessor(batch)
	}
}

func batchProcessor(batch []MoniObj) {
	fmt.Printf("接收到数据：%v \n", len(batch))
	result := make([]taosUtil.SubTableValue, 0)
	for _, obj := range batch {
		tags := []taosUtil.TagValue{
			{
				Name:  "device",
				Value: "device01",
			},
		}
		values := []taosUtil.RowValue{{
			Fields: []taosUtil.FieldValue{
				{
					Name:  "ts",
					Value: obj.ts,
				}, {
					Name:  "airStatus",
					Value: obj.airStatus,
				}, {
					Name:  "airElectric",
					Value: obj.airElectric,
				}, {
					Name:  "airGas",
					Value: obj.airGas,
				}, {
					Name:  "airWater",
					Value: obj.airWater,
				}, {
					Name:  "airEscapage",
					Value: obj.airEscapage,
				}, {
					Name:  "airFlow",
					Value: obj.airFlow,
				}, {
					Name:  "airPower",
					Value: obj.airPower,
				}, {
					Name:  "airExhaustPressure",
					Value: obj.airExhaustPressure,
				}, {
					Name:  "airGasDisplacement",
					Value: obj.airGasDisplacement,
				}, {
					Name:  "airMainPressure",
					Value: obj.airMainPressure,
				},
			},
		}}
		var subTableValue taosUtil.SubTableValue
		subTableValue.Name = "device01"
		subTableValue.SuperTable = "device"
		subTableValue.Tags = tags
		subTableValue.Values = values
		result = append(result, subTableValue)
	}
	_, err := taosUtil.InsertAutoCreateTable(result)
	if err != nil {
		fmt.Println("taos insert error :" + err.Error())
		panic(err.Error())
	}
}

func startMoni() {
	init := MoniObj{
		airStatus:          "Run",
		airElectric:        0,
		airGas:             0,
		airWater:           0,
		airEscapage:        0,
		airFlow:            0,
		airPower:           0,
		airExhaustPressure: 0,
		airGasDisplacement: 0,
		airMainPressure:    0,
		ts:                 time.Now().UnixNano() / 1e6,
	}
	for true {
		result := MoniObj{
			airStatus:          "Run",
			airElectric:        init.airElectric + getRandomFloat(0.04, 0.044),
			airGas:             init.airGas + getRandomFloat(0.3, 0.34),
			airWater:           init.airWater + getRandomFloat(0.008, 0.009),
			airEscapage:        init.airEscapage + getRandomFloat(0.002, 0.004),
			airFlow:            getRandomInt(0, 45),
			airPower:           getRandomInt(230, 258),
			airExhaustPressure: getRandomFloat(0.3, 0.4),
			airGasDisplacement: init.airGasDisplacement + getRandomFloat(0.03, 0.04),
			airMainPressure:    getRandomInt(0, 12),
			ts:                 time.Now().UnixNano() / 1e6,
		}
		init = result
		EventQueue <- init
		time.Sleep(1 * time.Second)
	}
}

func getRandomFloat(min, max float64) float64 {
	re := min + rand.Float64()*(max-min)
	sprintf := fmt.Sprintf("%.3f", re)
	float, _ := strconv.ParseFloat(sprintf, 64)
	return float
}

func getRandomInt(min, max int) int {
	return rand.Intn(max-min) + min
}
