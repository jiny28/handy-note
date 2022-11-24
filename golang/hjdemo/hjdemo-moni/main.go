package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var point = flag.Int("p", 10000, "point num")
var deviceNum = flag.Int("dn", 10, "device num")
var ms = flag.Int("ms", 10, "sleep ms")
var numSecond = flag.Int("cs", 5, "send m data")

func main() {
	flag.Parse()
	for i := 0; i < *deviceNum; i++ {
		device := "device" + strconv.Itoa(i)
		go startClient(device)
	}
	for {
		time.Sleep(time.Second * 100)
	}
}

func startClient(device string) {
	conn, err := net.Dial("tcp", "192.168.1.229:8888")
	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}
	fmt.Println("client connect 成功")
	defer conn.Close()
	num := *numSecond * 1000 / (*ms)
	fmt.Printf("共计发送%v次\n", num)
	now := time.Now()
	data := getDeviceData(device)
	bytes := []byte(data)
	for i := 0; i < num; i++ {
		_, err = conn.Write(bytes)
		if err != nil {
			fmt.Println("coon.write err=", err)
			break
		}
		time.Sleep(10 * time.Millisecond)
	}
	fmt.Printf("%v设备发送完毕,总计耗时%v\n", device, time.Since(now))
	for {
		time.Sleep(500 * time.Second)
	}
}

func getDeviceData(device string) string {
	value := "10.808"
	var builder strings.Builder
	builder.WriteString(device + ",")
	/*tsInt := time.Now().UnixNano() / 1e6
	builder.WriteString(strconv.FormatInt(tsInt, 10) + ",")*/
	for i := 0; i < *point; i++ {
		builder.WriteString(value)
		if i != *point-1 {
			builder.WriteString(",")
		}
	}
	builder.WriteString("E")
	//result := fmt.Sprintf("%x", builder.String())
	return builder.String()
}
