package main

import (
	"bi-demo/entity"
	"bi-demo/taosUtil"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var (
	port = 9090
)

type rep struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

func startWeb() {
	//127.0.0.1:9090/bi/real?device=air1,air4&code=airflow,airexhaustpressure,r_airelectric_day
	http.HandleFunc("/bi/real", realData)
	//127.0.0.1:9090/bi/getHisData?device=air1&code=airelectric,airescapage,r_airelectric_day&startTime=2022-11-23 23:00:00&endTime=2022-11-24 01:00:00
	http.HandleFunc("/bi/getHisData", getHisData)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		fmt.Println("http监听错误:" + err.Error())
		panic(err)
	}
}
func realData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	d, ok := r.Form["device"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未传设备."))
		return
	}
	c, ok := r.Form["code"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未传编码."))
		return
	}
	resultData := make(map[string]map[string]entity.RealValue, 0)
	devices := strings.Split(d[0], ",")
	codes := strings.Split(c[0], ",")
	for i := range devices {
		device := devices[i]
		for j := range codes {
			code := codes[j]
			key := device + "-" + code
			if strings.Contains(code, "r_") {
				key = device + "__" + code
			}
			load, ok := RealValue.Load(key)
			if !ok {
				continue
			}
			value := load.(entity.RealValue)
			m, ok := resultData[device]
			if !ok {
				m = make(map[string]entity.RealValue)
			}
			m[code] = value
			resultData[device] = m
		}
	}
	if len(resultData) == 0 {
		fmt.Fprintf(w, getErrorRes("未查询到数据!"))
		return
	}
	var resultReq = rep{Code: 200, Data: resultData}
	result, err := json.Marshal(resultReq)
	if err != nil {
		fmt.Fprintf(w, getErrorRes("序列化错误!"))
		return
	}
	fmt.Fprintf(w, string(result))
}

var backname = "biname1"

//CREATE STABLE biback (ts timestamp, code binary(64),hvalue binary(16000),device binary(64),starttime binary(64),endtime binary(64)) TAGS (backname binary(64));
func getHisData(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() // 解析参数，默认是不会解析的
	s, ok := r.Form["startTime"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择开始时间."))
		return
	}
	startTime := s[0]
	e, ok := r.Form["endTime"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择结束时间."))
		return
	}
	endTime := e[0]
	d, ok := r.Form["device"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择设备."))
		return
	}
	device := d[0]
	c, ok := r.Form["code"]
	if !ok {
		fmt.Fprintf(w, getErrorRes("未选择项."))
		return
	}
	codes := strings.Split(c[0], ",")
	// 筛选出已经存在副本的数据项
	filterCodes, backData := filterBack(codes, device, startTime, endTime)
	if len(filterCodes) > 0 {
		resultItems := make([]string, 0)
		itemCodes := make([]string, 0)
		for i := range filterCodes {
			code := filterCodes[i]
			if strings.Contains(code, "r_") {
				// result
				resultItems = append(resultItems, " item = '"+device+"__"+code+"' ")
			} else {
				// itemCode
				itemCodes = append(itemCodes, code)
			}
		}
		if len(resultItems) > 0 {
			// 有result
			join := "(" + strings.Join(resultItems, "or") + ")"
			resultSql := fmt.Sprintf("select * from kyjres where %v and ts >= '%v' and ts <= '%v' order by ts asc", join, startTime, endTime)
			query, _ := taosUtil.ExecuteQuery(resultSql, "")
			if query != nil && len(query) > 0 {
				for _, v := range query {
					key := strings.Split(v["item"], "__")[1]
					ts := v["ts"]
					hvalue := v["hvalue"]
					parseInt, _ := strconv.ParseInt(ts, 10, 64)
					time := time.Unix(0, parseInt*1e6).Format(layout)
					values, ok := backData[key]
					if !ok {
						values = make([]entity.RealValue, 0)
					}
					values = append(values, entity.RealValue{Date: time, Value: hvalue})
					backData[key] = values
				}
			}
		}
		if len(itemCodes) > 0 {
			// 有itemCodes
			join := strings.Join(itemCodes, ",")
			itemCodeSql := fmt.Sprintf("select ts,%v from kyj where ts >= '%v' and ts <= '%v' order by ts asc ", join, startTime, endTime)
			query, _ := taosUtil.ExecuteQuery(itemCodeSql, "")
			if query != nil && len(query) > 0 {
				for _, v := range query {
					ts := v["ts"]
					parseInt, _ := strconv.ParseInt(ts, 10, 64)
					time := time.Unix(0, parseInt*1e6).Format(layout)
					for vk, vv := range v {
						if vk == "ts" || vv == "" {
							continue
						}
						values, ok := backData[vk]
						if !ok {
							values = make([]entity.RealValue, 0)
						}
						values = append(values, entity.RealValue{Date: time, Value: vv})
						backData[vk] = values
					}
				}
			}
		}
	}
	if len(backData) > 0 {
		var resultReq = rep{Code: 200, Data: backData}
		result, err := json.Marshal(resultReq)
		if err != nil {
			fmt.Fprintf(w, getErrorRes("序列化错误!"))
			return
		}
		writeCacheTaos(backData, filterCodes, device, startTime, endTime)
		fmt.Fprintf(w, string(result))
		return
	} else {
		fmt.Fprintf(w, getErrorRes("未查询到数据!"))
		return
	}
}

func writeCacheTaos(data map[string][]entity.RealValue, codes []string, device string, startTime string, endTime string) {
	var subTableValue taosUtil.SubTableValue
	subTableValue.Name = backname
	subTableValue.SuperTable = "biback"
	tags := []taosUtil.TagValue{
		{
			Name:  "backname",
			Value: backname,
		},
	}
	subTableValue.Tags = tags
	rowValues := make([]taosUtil.RowValue, 0)
	for key := range data {
		if !strings.Contains(key, "r_") {
			continue
		}
		for i := range codes {
			code := codes[i]
			if code == key {
				now := time.Now().UnixNano() / 1e6
				values := data[key]
				marshal, _ := json.Marshal(&values)
				// 新数据
				fieldValues := []taosUtil.FieldValue{
					{
						Name:  "ts",
						Value: now,
					},
					{
						Name:  "code",
						Value: key,
					},
					{
						Name:  "hvalue",
						Value: string(marshal),
					},
					{
						Name:  "device",
						Value: device,
					},
					{
						Name:  "starttime",
						Value: startTime,
					},
					{
						Name:  "endtime",
						Value: endTime,
					},
				}
				rowValues = append(rowValues, taosUtil.RowValue{Fields: fieldValues})
			}
		}
	}
	if len(rowValues) == 0 {
		return
	}
	subTableValue.Values = rowValues
	_, err := taosUtil.InsertAutoCreateTable([]taosUtil.SubTableValue{subTableValue})
	if err != nil {
		fmt.Println("快照保存失败：" + err.Error())
		panic(err.Error())
	}
}

func filterBack(codes []string, device string, startTime string, endTime string) ([]string, map[string][]entity.RealValue) {
	sqlArray := make([]string, 0)
	for i := range codes {
		code := codes[i]
		sqlArray = append(sqlArray, " code = '"+code+"' ")
	}
	conditionSql := "(" + strings.Join(sqlArray, "or") + ")"
	sql := fmt.Sprintf("select code,hvalue from %v where %v and device = '%v' and starttime = '%v' and endtime = '%v'", backname, conditionSql, device, startTime, endTime)
	query, _ := taosUtil.ExecuteQuery(sql, "")
	result := make(map[string][]entity.RealValue)
	if query == nil || len(query) == 0 {
		return codes, result
	}
	for i := range query {
		m := query[i]
		mCode := m["code"]
		mValue := m["hvalue"]
		var realV []entity.RealValue
		err := json.Unmarshal([]byte(mValue), &realV)
		if err != nil {
			fmt.Println("副本数据序列化错误:" + err.Error())
			continue
		}
		result[mCode] = realV
	}
	resultCode := make([]string, 0)
	for i := range codes {
		c := codes[i]
		f := true
		for k := range result {
			if c == k {
				f = false
				break
			}
		}
		if f {
			resultCode = append(resultCode, c)
		}
	}
	return resultCode, result
}

func getErrorRes(msg string) string {
	var errorReq = rep{Code: 0, Data: msg}
	marshal, _ := json.Marshal(errorReq)
	return string(marshal)
}
