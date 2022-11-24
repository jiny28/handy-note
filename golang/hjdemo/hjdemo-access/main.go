package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"sync"
	"time"
)

var buffer = 2048

func main() {
	l, err := net.Listen("tcp", ":8888")
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer l.Close()
	runDeviceQueue()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
		fmt.Println("建立一个客户端连接.")
		go handleConn(c)
	}
}

var dealLock sync.Mutex
var dealSumTime time.Duration
var mqttSumTime time.Duration

func handleConn(c net.Conn) {
	var builder strings.Builder
	nowTime := time.Now()
	defer c.Close()
	for {
		var byt = make([]byte, buffer)
		n, err := c.Read(byt)
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端断开连接.")
			} else {
				fmt.Println("conn read error:", err)
			}
			return
		}
		hexData := hex.EncodeToString(byt[0:n])
		index := strings.Index(hexData, "45")
		if index == -1 {
			// 没到该设备的结尾
			builder.WriteString(hexData)
		} else {
			builder.WriteString(hexData[:index+2])
			deviceData := builder.String()
			builder.Reset()
			builder.WriteString(hexData[index+2:])
			deviceHexByte, _ := hex.DecodeString(deviceData)
			endData := string(deviceHexByte)
			resultData := endData[:len(endData)-1]
			//fmt.Println("读取到设备数据:" + resultData)
			//fmt.Printf("解析一台设备数据所耗时:%v\n", time.Since(nowTime))
			datas := splitData(resultData, 2500)
			var mqttS time.Duration
			for i := range datas {
				splitStrData := datas[i]
				EventQueue <- splitStrData
				split := strings.Split(splitStrData, ",")
				device := split[0]
				nowMs := time.Now().UnixNano() / 1e6
				mqttNow := time.Now()
				mqttError := mqttConnection.PublishMsg("exdevice/"+device+"/"+strconv.FormatInt(nowMs, 10), 0, false, splitStrData)
				since := time.Since(mqttNow)
				mqttS += since
				if mqttError != nil {
					fmt.Printf("mqtt转发错误device:%v:%v\n", device, mqttError.Error())
				}
			}
			dealLock.Lock()
			mqttSumTime += mqttS
			dealTime := time.Since(nowTime) - mqttS
			dealSumTime += dealTime
			dealLock.Unlock()
			nowTime = time.Now()
		}
	}
}

func splitData(data string, num int64) []string {
	split := strings.Split(data, ",")
	device := split[0]
	datas := split[1:]
	lens := int64(len(datas))
	if lens <= num {
		return []string{data}
	}
	//获取应该数组分割为多少份
	var splitResult = make([]string, 0)
	var quantity int64
	if lens%num == 0 {
		quantity = lens / num
	} else {
		quantity = (lens / num) + 1
	}
	//声明分割数组的截止下标
	var start, end, i int64
	for i = 1; i <= quantity; i++ {
		end = i * num
		if i != quantity {
			value := device + "," + strings.Join(datas[start:end], ",")
			splitResult = append(splitResult, value)

		} else {
			value := device + "," + strings.Join(datas[start:], ",")
			splitResult = append(splitResult, value)
		}
		start = i * num
	}
	return splitResult
}
