package main

import (
	"encoding/json"
	"fmt"
	"iot-demo/taosUtil"
	"net/http"
	"strconv"
)

var (
	port = 8080
)

type rep struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func startWeb() {
	http.HandleFunc("/goiot/real", realData)
	http.HandleFunc("/goiot/getData", getData)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("http监听错误:" + err.Error())
		panic(err)
	}
}

func getData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	startTime, ok := r.Form["startTime"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择开始时间."))
		return
	}
	endTime, ok := r.Form["endTime"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择结束时间."))
		return
	}
	items, ok := r.Form["items"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择采集项."))
		return
	}
	sql := fmt.Sprintf("select ts,%v from device where ts >= '%v' and ts <= '%v' order by ts asc", items[0], startTime[0], endTime[0])

	query, _ := taosUtil.ExecuteQuery(sql, "")
	if query == nil || len(query) == 0 {
		fmt.Fprintf(w, getErrorRes("未查询到数据!"))
		return
	}
	var convertData = make(map[string][]interface{})
	for _, v := range query {
		for vk, vv := range v {
			value, ok := convertData[vk]
			if !ok {
				value = make([]interface{}, 0)
			}
			value = append(value, vv)
			convertData[vk] = value
		}
	}
	var resultReq = rep{Code: 200, Data: convertData}
	result, err := json.Marshal(resultReq)
	if err != nil {
		fmt.Fprintf(w, getErrorRes("序列化错误!"))
		return
	}
	fmt.Fprintf(w, string(result))
}

func realData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	device, ok := r.Form["device"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择开始时间."))
		return
	}
	sql := fmt.Sprintf("select last_row(*) from %v", device[0])
	query, _ := taosUtil.ExecuteQuery(sql, "")
	if query == nil || len(query) == 0 {
		fmt.Fprintf(w, getErrorRes("未查询到数据!"))
		return
	}
	var resultReq = rep{Code: 200, Data: query[0]}
	result, err := json.Marshal(resultReq)
	if err != nil {
		fmt.Fprintf(w, getErrorRes("序列化错误!"))
		return
	}
	fmt.Fprintf(w, string(result))
}

func getErrorRes(msg string) string {
	var errorReq = rep{Code: 0, Data: msg}
	marshal, _ := json.Marshal(errorReq)
	return string(marshal)
}
