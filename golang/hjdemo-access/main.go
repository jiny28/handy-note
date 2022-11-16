package main

import (
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"strings"
)

var buffer = 4000

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

func handleConn(c net.Conn) {
	ip := strings.Split(c.RemoteAddr().String(), ":")[0]
	fmt.Println("开始处理ip:" + ip)
	var builder strings.Builder
	defer fmt.Println("线程结束：" + ip)
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
			EventQueue <- resultData
		}
	}
}
