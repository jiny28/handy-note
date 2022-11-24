package main

import (
	"bi-demo/taosUtil"
	"fmt"
)

var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "h",
}

func main() {
	// 开启实时的数据
	startSub()
	startWeb()
}

func init() {
	taosUtil.Connection(taosInfo)
}

func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
