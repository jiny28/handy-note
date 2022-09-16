package task

import (
	"clear-data/mysqlUtil"
	"clear-data/taosUtil"
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"strconv"
	"strings"
	"time"
)

var mysqlInfo = mysqlUtil.MysqlInfo{
	UserName: "root",
	Password: "123456",
	Ip:       "10.88.0.14",
	Port:     33066,
	Db:       "hlhz_go_test",
}
var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "hlhz1",
}

var globalIndex, globalTotal int64
var myScadaTypes []string
var startTime int64
var itemInfos map[string]int
var itemMapScadaTypes map[string]string
var itemMapScadaCodes map[string]string
var valueMap map[string]map[string]string
var lastTimeMap map[string]string
var lastValueMap map[string]map[string]string

var typeMap map[string]string
var dataRate map[string]string
var fullValue map[string]string

var sourceDbName = "hlhz"
var lastValueMapName = "lastValueMap"
var interval int64 = 1000 * 60 * 1
var currentTime int64 = 1000 * 60 * 1
var layout = "2006-01-02 15:04:05"

type Items struct {
	item         string //原始数据项表名
	time         string //采集时间
	value        string //采集值
	fullInterval int64  //补数据的时长
	currentTime  int64  //距离当前时间的时长
}

func init() {
	_, error := taosUtil.Connection(taosInfo)
	if error != nil {
		fmt.Println("taos  连接错误:" + error.Error())
		panic(error)
	}
	_, error = mysqlUtil.Connection(mysqlInfo)
	if error != nil {
		fmt.Println("mysql 连接错误:" + error.Error())
		panic(error)
	}
}

func ClearData(cxt context.Context, param *xxl.RunReq) (msg string) {
	index := param.BroadcastIndex
	total := param.BroadcastTotal
	params := param.ExecutorParams
	// 限制启动的执行器个数
	if params != "" {
		parseInt, err := strconv.ParseInt(params, 10, 64)
		if err != nil {
			e := "入参错误:" + err.Error()
			fmt.Println(e)
			return e
		}
		if parseInt < total {
			total = parseInt
		}
	}
	if index >= total {
		return "超出设定执行器总数,无需执行"
	}
	routInit(index, total)
	if myScadaTypes == nil || len(myScadaTypes) == 0 {
		return "无采集器类型"
	}
	if startTime == 0 {
		return "无开始时间"
	}
	endTime := startTime + interval
	nowTime := time.Now().UnixNano() / 1e6
	if endTime > nowTime-currentTime {
		endTime = nowTime - currentTime
	}
	fmt.Sprintf("清洗数据的开始时间为:%v,结束时间为%v", time.Unix(startTime, 0).Format(layout), time.Unix(endTime, 0).Format(layout))
	if itemInfos == nil {
		return "itemInfos 没数据"
	}
	data := getData(myScadaTypes, startTime, endTime)
	// 深拷贝
	copyItemInfos := make(map[string]int, len(itemInfos))
	for k, v := range itemInfos {
		copyItemInfos[k] = v
	}
	items := formatData(data, copyItemInfos)
	clearInit(index, total)
	if typeMap == nil || len(typeMap) == 0 {
		return "typeMap is empty"
	}

	for _, item := range items {
		dealData(item)
	}
	writeTaos(itemMapScadaTypes, itemMapScadaCodes)
	setVarToTao()
	startTime = endTime
	globalIndex = index
	globalTotal = total
	return
}

func setVarToTao() {
	if lastValueMap == nil || len(lastValueMap) == 0 {
		return
	}
	insertTime := "2022-08-15 00:00:00"
	parseTime, _ := time.ParseInLocation(layout, insertTime, time.Local)
	insertTimeLong := parseTime.UnixNano() / 1e6
	var result = make([]taosUtil.SubTableValue, 0)
	for tableName, v := range lastValueMap {
		time := v["time"]
		value := v["value"]
		var subTableValue taosUtil.SubTableValue
		subTableValue.Name = tableName
		subTableValue.SuperTable = lastValueMapName
		tags := []taosUtil.TagValue{
			{
				Name:  "location",
				Value: tableName,
			},
		}
		subTableValue.Tags = tags
		rowValues := []taosUtil.RowValue{
			{
				Fields: []taosUtil.FieldValue{
					{
						Name:  "ts",
						Value: insertTimeLong,
					},
					{
						Name:  "time",
						Value: time,
					},
					{
						Name:  "value",
						Value: value,
					},
				},
			},
		}
		subTableValue.Values = rowValues
		result = append(result, subTableValue)
	}
	_, e := taosUtil.InsertAutoCreateTable(result)
	if e != nil {
		fmt.Printf("taos insert error:" + e.Error())
	}
}

func writeTaos(types map[string]string, codes map[string]string) {
	if valueMap == nil || len(valueMap) == 0 {
		return
	}
	// valueMap - k:device_itemCode v:map<time,value>
	// valueMap -> map[string]map[int64]map[string]interface{}
	// <device_scada,<longTs,<fieldName,fieldValue>>>
	taosData := make(map[string]map[int64]map[string]interface{})
	for k, v := range valueMap {
		split := strings.Split(k, "_")
		device := split[0]
		code := split[1]
		stable, ok := types[code]
		if !ok || stable == "" {
			fmt.Println(k + "没有对应的超级表")
			continue
		}
		device = device + "_" + stable
		vt, ok := taosData[device]
		tmap := make(map[int64]map[string]interface{})
		if ok {
			tmap = vt
			for vk, vv := range v {
				vvSplit := strings.Split(vv, "_")
				value := vvSplit[0]
				vkTime, _ := time.ParseInLocation(layout, vk, time.Local)
				timeLong := vkTime.UnixNano() / 1e6
				vvt, ok := tmap[timeLong]
				var vmap = make(map[string]interface{})
				if ok {
					vmap = vvt
				}
				vmap[code] = value
				tmap[timeLong] = vmap
			}
		} else {
			for vk, vv := range v {
				vvSplit := strings.Split(vv, "_")
				value := vvSplit[0]
				vkTime, _ := time.ParseInLocation(layout, vk, time.Local)
				timeLong := vkTime.UnixNano() / 1e6
				var vmap = make(map[string]interface{})
				vmap[code] = value
				tmap[timeLong] = vmap
			}
		}
		taosData[device] = tmap
	}
	subTableValue := gerSubTableValue(taosData)
	_, e := taosUtil.InsertAutoCreateTable(subTableValue)
	if e != nil {
		fmt.Printf("taos insert error:" + e.Error())
	}
}

func gerSubTableValue(data map[string]map[int64]map[string]interface{}) []taosUtil.SubTableValue {
	var result = make([]taosUtil.SubTableValue, 0)
	for device, v := range data {
		var subTableValue taosUtil.SubTableValue
		subTableValue.Name = device
		split := strings.Split(device, "_")
		subTableValue.SuperTable = split[1]
		tags := []taosUtil.TagValue{
			{
				Name:  "device",
				Value: split[0],
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

func dealData(item Items) {
	tableName := item.item
	time := item.time
	value := item.value
	clearType, ok := typeMap[tableName]
	if !ok {
		return
	}
	if time == "" {
		// 补数据
		fullNull(tableName, clearType, item.fullInterval, item.currentTime)
		return
	}
	if value == "" {
		value = "无"
	}
	switch clearType {
	case "ct01":
		dealAdd(time, value, tableName, clearType)
	case "ct02":
		dealStatus(time, value, tableName, clearType)
	case "ct03":
		dealStatic(time, value, tableName)
	}
}

func dealStatic(time string, value string, tableName string) {
	vMap, ok := valueMap[tableName]
	if !ok {
		vMap = make(map[string]string)
	}
	vMap[time] = value
	lMap := make(map[string]string)
	lMap["time"] = time
	lMap["value"] = value
	valueMap[tableName] = vMap
	lastValueMap[tableName] = lMap
}

func dealStatus(time string, value string, tableName string, clearType string) {
	lMap, ok := lastValueMap[tableName]
	if !ok {
		lMap = make(map[string]string)
	}
	vMap, ok := valueMap[tableName]
	if !ok {
		vMap = make(map[string]string)
	}
	lastTime, tOk := lMap["time"]
	lastValue, vOk := lMap["value"]
	if !tOk || !vOk {
		lastTime = time
		lastValue = value
		vMap[lastTime] = lastValue
	}
	k := intervalTime(lastTime, time)
	if k < 0 {
		return
	}
	last, ok := lastTimeMap[tableName]
	if !ok {
		last = time
	}
	j := intervalTime(last, time)
	rate := getValue(dataRate, tableName, "60")
	rateI, _ := strconv.Atoi(rate)
	if value != lastValue {
		if j > int64(rateI) {
			lastValue = getValue(fullValue, tableName, lastValue)
		}
		fullData(lastTime, time, lastValue, value, rateI, "1", vMap, clearType)
		vMap[time] = value
		lastTime = time
		lastValue = value
	} else {
		addTime := fullData(lastTime, time, lastValue, value, rateI, "1", vMap, clearType)
		lastTime = addTime
	}
	lastTimeMap[tableName] = time
	lMap["time"] = lastTime
	lMap["value"] = lastValue
	//存放 value_map
	valueMap[tableName] = vMap
	lastValueMap[tableName] = lMap
}

func dealAdd(time string, value string, tableName string, clearType string) {
	lMap, ok := lastValueMap[tableName]
	if !ok {
		lMap = make(map[string]string)
	}
	vMap, ok := valueMap[tableName]
	if !ok {
		vMap = make(map[string]string)
	}
	lastTime, tOk := lMap["time"]
	lastValue, vOk := lMap["value"]
	if !tOk || !vOk {
		lastTime = time
		lastValue = value
		vMap[lastTime] = lastValue
	}
	k := intervalTime(lastTime, time)
	if k < 0 {
		return
	}
	last, ok := lastTimeMap[tableName]
	if !ok {
		last = time
	}
	j := intervalTime(last, time)
	rate := getValue(dataRate, tableName, "60")
	rateI, _ := strconv.Atoi(rate)
	if k > int64(rateI) {
		if j > int64(rateI) {
			lastValue = getValue(fullValue, tableName, lastValue)
		}
		// 截取60秒保存该值
		addTime := fullData(lastTime, time, lastValue, value, rateI, "1", vMap, clearType)
		lastTime = addTime
		lastValue = value
	} else if k == int64(rateI) {
		vMap[time] = value
		lastTime = time
		lastValue = value
	} else {
		lastValue = value

	}
	lastTimeMap[tableName] = time
	lMap["time"] = lastTime
	lMap["value"] = lastValue
	//存放 value_map
	valueMap[tableName] = vMap
	lastValueMap[tableName] = lMap
}
func intervalTime(startTime string, endTime string) int64 {
	start, _ := time.ParseInLocation(layout, startTime, time.Local)
	end, _ := time.ParseInLocation(layout, endTime, time.Local)
	return (end.UnixNano()/1e6 - start.UnixNano()/1e6) / 1000
}

func fullNull(tableName string, clearType string, fullInterval int64, currentTime int64) {
	lMap, ok := lastValueMap[tableName]
	if !ok {
		lMap = make(map[string]string)
	}
	vMap, ok := valueMap[tableName]
	if !ok {
		vMap = make(map[string]string)
	}
	lastTime, tOk := lMap["time"]
	lastValue, vOk := lMap["value"]
	if !tOk || !vOk {
		return
	}
	addTime := addSecondTime(lastTime, fullInterval/1000)
	now := time.Unix(time.Now().UnixNano()/1e6-currentTime, 0).Format(layout)
	b := compareTime(now, addTime)
	if !b {
		addTime = now
	}
	rate := getValue(dataRate, tableName, "60")
	rateI, _ := strconv.Atoi(rate)
	lastValue = getValue(fullValue, tableName, lastValue)
	lastTime = fullData(lastTime, addTime, lastValue, lastValue, rateI, "2", vMap, clearType)
	lMap["time"] = lastTime
	lMap["value"] = lastValue
	lastValueMap[tableName] = lMap
	valueMap[tableName] = vMap

}

func fullData(start string, end string, value string, value1 string, addTime int, flag string, valueMap map[string]string, clearType string) string {
	flags := true
	i := 0
	for flags {
		start = addSecondTime(start, int64(addTime))
		b := compareTime(start, end)
		if b {
			flags = false
		} else {
			if start == end {
				if "2" == flag {
					valueMap[start] = value + "_" + flag
				} else {
					valueMap[start] = value1
				}
			} else {
				valueMap[start] = value + "_" + flag
			}
			i++
		}
	}
	//最后一条时间
	start = addSecondTime(start, int64(addTime*-1))
	if i == 1 && "1" == flag && "0" == clearType {
		valueMap[start] = value
	}
	return start
}

func getValue(vMap map[string]string, tableName string, defaultValue string) string {
	value, ok := vMap[tableName]
	if !ok || value == "" {
		value = defaultValue
	}
	return value
}

func compareTime(startTime string, endTime string) bool {
	sTime, _ := time.ParseInLocation(layout, startTime, time.Local)
	eTime, _ := time.ParseInLocation(layout, endTime, time.Local)
	if sTime.After(eTime) {
		return true
	}
	return false
}

func addSecondTime(lastTime string, i int64) string {
	lTime, _ := time.ParseInLocation(layout, lastTime, time.Local)
	return lTime.Add(time.Second * time.Duration(i)).Format(layout)
}

func clearInit(index int64, total int64) {
	valueMap = make(map[string]map[string]string)
	if index != globalIndex || total != globalTotal {
		initLastValueMap()
	}
}

func initLastValueMap() {
	sql := "select * from " + lastValueMapName + ""
	query, _ := taosUtil.ExecuteQuery(sql, "")
	if query == nil || len(query) == 0 {
		return
	}
	lastTimeMap = make(map[string]string)
	lastValueMap = make(map[string]map[string]string)
	for _, v := range query {
		di := v["location"]
		time := v["time"]
		lastTimeMap[di] = time
		lastValueMap[di] = v
	}
}

func routInit(index int64, total int64) {
	if index != globalIndex || total != globalTotal {
		// init
		myScadaTypes = nil
		startTime = 0
		itemInfos, itemMapScadaTypes, itemMapScadaCodes, fullValue, dataRate, typeMap = nil, nil, nil, nil, nil, nil
		// 获取所有采集器类型基础数据，执行器按照采集器类型平分数量
		scadaTypes := getScadaTypeInfo()
		if scadaTypes == nil {
			fmt.Println("scadaType 未配置")
			return
		}
		splitScadaTypes := splitArray(scadaTypes, total)
		myScadaTypes = splitScadaTypes[index]
		fmt.Sprintf("下标为%v的执行器正在执行以下表名%v", index, myScadaTypes)
		startTime = getMinTime(myScadaTypes)
		itemInfos, itemMapScadaTypes, itemMapScadaCodes, fullValue, dataRate, typeMap = getItemInfos(myScadaTypes)
	}
}

func formatData(data []map[string]string, itemInfos map[string]int) []Items {
	result := make([]Items, 0)
	for _, v := range data {
		ts := v["ts"]
		device := v["device"]
		timeL, _ := strconv.ParseInt(ts, 10, 64)
		time := time.Unix(timeL, 0).Format(layout)
		delete(v, "ts")
		delete(v, "device")
		for k, value := range v {
			items := Items{time: time, item: device + "_" + k, value: value}
			delete(itemInfos, device+"_"+k)
			result = append(result, items)
		}
	}
	if len(itemInfos) > 0 {
		for k, _ := range itemInfos {
			items := Items{item: k, fullInterval: interval, currentTime: currentTime}
			result = append(result, items)
		}
	}
	return result
}

func getItemInfos(types []string) (map[string]int, map[string]string, map[string]string, map[string]string, map[string]string, map[string]string) {
	sql := "select i.c_scada_code,i.c_code,i.c_field_type,s.c_scada_type,i.c_device,i.c_clear_interval,i.c_clear_type,i.c_fill_data from t_item_info i " +
		"left join t_scada_info s on i.c_scada_code=s.c_scada_code where s.c_scada_type in("
	for i := range types {
		sql = sql + "'" + types[i] + "',"
	}
	sql = sql[0:len(sql)-1] + ")"
	all, _ := mysqlUtil.GetAll(sql)
	if all == nil || len(all) == 0 {
		return nil, nil, nil, nil, nil, nil
	}
	itemInfos := make(map[string]int)
	itemMapScadaTypes := make(map[string]string)
	itemMapScadaCodes := make(map[string]string)
	fullValue, dataRate, typeMap := make(map[string]string), make(map[string]string), make(map[string]string)
	for _, v := range all {
		itemCode := v["c_code"]
		device := v["c_device"]
		scadaType := v["c_scada_type"]
		scadaCode := v["c_scada_code"]
		fValue := v["c_fill_data"]
		rate := v["c_clear_interval"]
		cType := v["c_clear_type"]
		itemMapScadaTypes[itemCode] = scadaType
		itemMapScadaCodes[device+"_"+itemCode] = scadaCode
		itemInfos[device+"_"+itemCode] = 0
		fullValue[device+"_"+itemCode] = fValue
		dataRate[device+"_"+itemCode] = rate
		typeMap[device+"_"+itemCode] = cType
	}
	return itemInfos, itemMapScadaTypes, itemMapScadaCodes, fullValue, dataRate, typeMap
}

func getData(types []string, startTime int64, endTime int64) []map[string]string {
	result := make([]map[string]string, 0)
	for _, v := range types {
		sql := "select * from " + v + " where ts>=" + strconv.FormatInt(startTime, 10) + " and ts<" + strconv.FormatInt(endTime, 10)
		query, _ := taosUtil.ExecuteQuery(sql, sourceDbName)
		if query != nil {
			result = append(result, query...)
		}
	}
	return result
}

func getScadaTypeInfo() []string {
	sql := "select distinct c_scada_type from t_scada_info where c_flag='1' ORDER BY c_scada_type"
	all, _ := mysqlUtil.GetAll(sql)
	if len(all) == 0 {
		return nil
	}
	result := make([]string, len(all))
	for i, v := range all {
		result[i] = v["c_scada_type"]
	}
	return result
}

func splitArray(source []string, n int64) [][]string {
	sLen := int64(len(source))
	var result = make([][]string, 0)
	remainder := sLen % n
	number := sLen / n
	var offset int64 = 0 //偏移量
	var i int64 = 0
	for ; i < n; i++ {
		var value []string
		if remainder > 0 {
			value = source[i*number+offset : (i+1)*number+offset+1]
			remainder--
			offset++
		} else {
			value = source[i*number+offset : (i+1)*number+offset]
		}
		result = append(result, value)
	}
	return result
}

func getMinTime(myScadaTypes []string) int64 {
	times := make([]int64, 0)
	for _, v := range myScadaTypes {
		query, _ := taosUtil.ExecuteQuery("select last(*) from "+v+" group by device", "")
		if query != nil {
			ts := query[0]["ts"]
			tsInt, _ := strconv.ParseInt(ts, 10, 64)
			times = append(times, tsInt)
		}
	}
	if len(times) == 0 {
		for _, v := range myScadaTypes {
			query, _ := taosUtil.ExecuteQuery("select first(*) from "+v+" group by device", sourceDbName)
			if query != nil {
				ts := query[0]["ts"]
				tsInt, _ := strconv.ParseInt(ts, 10, 64)
				times = append(times, tsInt)
			}
		}
	}
	if len(times) == 0 {
		return 0
	}
	item := times[0]
	for _, v := range times {
		if item >= v {
			item = v
		}
	}
	return item
}
