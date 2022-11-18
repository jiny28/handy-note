package main

import (
	"flag"
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"
)

var point = flag.Int("p", 20, "point num")
var deviceNum = flag.Int("dn", 1, "device num")
var ms = flag.Int("ms", 100, "sleep ms")

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

	for {
		data := getDeviceData(device)
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("coon.write err=", err)
			break
		}
		time.Sleep(time.Duration(*ms) * time.Millisecond)
	}
}

func getDeviceData(device string) string {
	value := "10.808"
	var builder strings.Builder
	builder.WriteString(device + ",")
	tsInt := time.Now().UnixNano() / 1e6
	builder.WriteString(strconv.FormatInt(tsInt, 10) + ",")
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
