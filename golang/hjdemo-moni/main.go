package main

import (
	"flag"
	"fmt"
	"net"
	"strings"
	"time"
)

var point = flag.Int("p", 20, "point num")
var device = flag.String("d", "device0", "device name.")

func main() {
	flag.Parse()
	conn, err := net.Dial("tcp", "127.0.0.1:8888")
	if err != nil {
		fmt.Println("client dial err=", err)
		return
	}
	fmt.Println("client connect 成功")
	defer conn.Close()

	for {
		data := getDeviceData()
		_, err = conn.Write([]byte(data))
		if err != nil {
			fmt.Println("coon.write err=", err)
			break
		}
		time.Sleep(time.Millisecond * 100)
	}
}

func getDeviceData() string {
	value := "10.808"
	var builder strings.Builder
	builder.WriteString(*device + ",")
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
