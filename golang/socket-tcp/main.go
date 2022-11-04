package main

import (
	"database/sql"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"net"
	"socket-tcp/mysqlUtil"
	"socket-tcp/taosUtil"
	"strconv"
	"strings"
	"time"
)

var buffer = 2048
var port = ":60382"
var group = 500 // 一组队列的数量
var mysqlInfo = mysqlUtil.MysqlInfo{
	UserName: "root",
	Password: "123456",
	Ip:       "10.88.0.14",
	Port:     33066,
	Db:       "data",
}
var taosInfo = taosUtil.TaosInfo{
	HostName:   "taos-server",
	ServerPort: 6030,
	User:       "root",
	Password:   "taosdata",
	DbName:     "hlhz",
}

var db *sql.DB
var deviceIpMapping map[string]string
var deviceCount map[string]int
var itemMapping map[string][]map[string]string
var resultData map[string][]string

func init() {
	db = initMysql()
	taosUtil.Connection(taosInfo)
	list, err := mysqlUtil.GetAll("select * from t_device")
	checkErr(err, "get info sql")
	deviceIpMapping = make(map[string]string)
	deviceCount = make(map[string]int)
	itemMapping = make(map[string][]map[string]string)
	resultData = make(map[string][]string)
	for _, v := range list {
		device := v["c_device"]
		ip := v["c_ip"]
		deviceIpMapping[ip] = device
	}
	itemList, err := mysqlUtil.GetAll("select * from t_item_info  order by c_order+0 ")
	checkErr(err, "get info sql")
	for _, item := range itemList {
		item_code := item["item_code"]
		device_code := item["device_code"]
		data_type := item["data_type"]
		byte_length := item["byte_length"]
		m := make(map[string]string)
		m["item_code"] = item_code
		m["data_type"] = data_type
		m["byte_length"] = byte_length
		arrList, ok := itemMapping[device_code]
		if !ok {
			arrList = make([]map[string]string, 0)
		}
		arrList = append(arrList, m)
		itemMapping[device_code] = arrList
	}
}
func initMysql() *sql.DB {
	db, error := mysqlUtil.Connection(mysqlInfo)
	checkErr(error, "初始化MySql错误")
	fmt.Println(" mysql connection success")
	return db
}

var pData = flag.Int("p", 1, "print data log")

func main() {
	flag.Parse()
	l, err := net.Listen("tcp", port)
	if err != nil {
		fmt.Println("listen error:", err)
		return
	}
	defer l.Close()
	defer db.Close()
	RunQueue()
	for {
		c, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			break
		}
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
		if buffer == n { // 缓冲读满  直接转16进制
			builder.WriteString(hex.EncodeToString(byt))
		} else { //缓冲未占满 （最后一条数据）  需要线对字节进行切割
			builder.WriteString(hex.EncodeToString(byt[0:n]))
		}
		dealData(&builder, ip)
		//fmt.Printf("read %d bytes, content is %s\n", n, string(byt[:n]))
	}
}

func dealData(data *strings.Builder, ip string) {
	flag := true
	for flag {
		//头类型
		headerType := jugdeHeader(data.String())
		if headerType == 16 {
			// 获取数据类型
			dataType := data.String()[12:16]
			eventDataLen := hexStringToAlgorism(lowHighChange(data.String()[8:12])) * 2 //数据长度（头信息+数据体）
			if data.Len() < eventDataLen {
				//数据不全，结束
				flag = false
				continue
			}
			oneData := data.String()[16*2 : eventDataLen]
			if dataType == "1600" || dataType == "1700" { //事件数据
				readHexHaveHead(oneData, deviceIpMapping[ip])
			} else {
				readHexHaveHeadStatus(oneData, deviceIpMapping[ip])
			}
			rem := data.String()[eventDataLen:]
			data.Reset()
			data.WriteString(rem)
			if data.Len() < eventDataLen {
				//数据不全，结束
				flag = false
			}
		} else { //12字节头
			// 获取数据类型
			dataType := data.String()[0:4]                                                    //数据类型
			eventDataLen := (hexStringToAlgorism(lowHighChange(data.String()[4:8])) + 12) * 2 //数据长度（头信息+数据体）
			if data.Len() < eventDataLen {
				//数据不全，结束
				flag = false
				continue
			}
			//从请求头结束开始切
			oneData := data.String()[12*2 : eventDataLen]
			if dataType == "1600" || dataType == "1700" { //事件数据
				readHexHaveHead(oneData, deviceIpMapping[ip])
			} else {
				readHexHaveHeadStatus(oneData, deviceIpMapping[ip])
			}
			rem := data.String()[eventDataLen:]
			data.Reset()
			data.WriteString(rem)
			if data.Len() < eventDataLen {
				//数据不全，结束
				flag = false
			}
		}
	}

}

func readHexHaveHeadStatus(data string, device string) {
	itemInfo := itemMapping[device]
	startIndex := 0 //byte截取下标
	var sb strings.Builder
	us := getUs() //获取微妙

	sb.WriteString(device + ",")                    //设备编码保存到头部
	sb.WriteString(strconv.FormatInt(us, 10) + ",") //时间保存到头部
	for _, itemMap := range itemInfo {

		//            String itemCode = itemMap.get("item_code");//采集项编号
		dataType := itemMap["data_type"] //字节类型
		atoi, _ := strconv.Atoi(itemMap["byte_length"])
		length := atoi * 2 //字符长度  字节长度*2
		//原始byte截取下标
		substring := data[startIndex : startIndex+length]
		startIndex += length //下标偏移
		// （大小端）高低位转换
		dataLHex := lowHighChange(substring)
		var value string
		if dataType == "F32" {
			value = fmt.Sprintf("%v", toFloat(dataLHex))
		} else {
			value = fmt.Sprintf("%v", toInt(dataLHex))
		}
		sb.WriteString(value)
		sb.WriteString(",")
	}
	if *pData == 1 {
		fmt.Println("打印状态数据：" + sb.String())
	}
}

func readHexHaveHead(data string, device string) {
	//插入条数  上一次统计
	lastcount, ok := deviceCount[device]
	if !ok {
		lastcount = 0
		deviceCount[device] = lastcount
	}
	nowCount := lastcount + 1
	itemInfo := itemMapping[device]
	startIndex := 0 //byte截取下标
	var sb strings.Builder
	us := getUs() //获取微妙
	sb.WriteString(strconv.FormatInt(us, 10) + ",")
	for _, itemMap := range itemInfo {
		dataType := itemMap["data_type"] //字节类型
		atoi, _ := strconv.Atoi(itemMap["byte_length"])
		length := atoi * 2 //字符长度  字节长度*2
		//原始byte截取下标
		substring := data[startIndex : startIndex+length]
		startIndex += length //下标偏移
		// （大小端）高低位转换
		dataLHex := lowHighChange(substring)
		var value string
		if dataType == "F32" {
			value = fmt.Sprintf("%v", toFloat(dataLHex))
		} else {
			value = fmt.Sprintf("%v", toInt(dataLHex))
		}
		sb.WriteString(value)
		sb.WriteString(",")
	}
	strData := sb.String()[0 : sb.Len()-1]
	if *pData == 1 {
		fmt.Println("打印事件数据:" + sb.String())
	}
	listStr, ok := resultData[device]
	if !ok {
		listStr = make([]string, 0)
	}
	listStr = append(listStr, strData)
	resultData[device] = listStr
	if group == nowCount {
		deviceCount[device] = 0
		poll := resultData[device]
		m := make(map[string][]string)
		m[device] = poll
		EventQueue <- m
		delete(resultData, device)
	} else {
		deviceCount[device] = nowCount
	}
}

func toInt(str string) int {
	i, err := strconv.ParseInt(str, 16, 32)
	checkErr(err, "转换值错误")
	return int(i)
}

func toFloat(str string) float64 {
	return 0
}

func getUs() int64 {
	return time.Now().UnixNano() / 1e3
}

func jugdeHeader(resHexStr string) int {
	headerType := 12
	headHex := resHexStr[0:8]
	if headHex == "44332211" {
		headerType = 16
	}
	return headerType
}

func lowHighChange(str string) string {
	chars := strings.Split(str, "")
	var res string
	for i := len(chars) - 1; i > 0; i-- {
		if i%2 == 1 { //第偶数个
			data1 := chars[i-1]
			data2 := chars[i]
			res += data1 + data2
		}
	}
	return res
}

func hexStringToAlgorism(hex string) int {
	parseUint, err := strconv.ParseUint(hex, 16, 32)
	if err != nil {
		fmt.Println(err.Error())
	}
	return int(parseUint)
}
func checkErr(err error, prompt string) {
	if err != nil {
		fmt.Printf("error: %s\n", prompt)
		panic(err.Error())
	}
}
