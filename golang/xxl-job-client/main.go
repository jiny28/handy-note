package main

import (
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"log"
	"xxl-job-client/task"
)

func main() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr("http://10.88.0.14:28888/hlhz_task_admin"),
		xxl.AccessToken(""),             //请求令牌(默认为空)
		xxl.ExecutorIp("10.88.0.141"),   //可自动获取
		xxl.ExecutorPort("9996"),        //默认9999（非必填）
		xxl.RegistryKey("golang-clear"), //执行器名称
		xxl.SetLogger(&logger{}),        //自定义日志
	)
	exec.Init()
	//设置日志查看handler
	exec.LogHandler(func(req *xxl.LogReq) *xxl.LogRes {
		return &xxl.LogRes{Code: 200, Msg: "", Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "这个是自定义日志handler",
			IsEnd:       true,
		}}
	})
	//注册任务handler
	exec.RegTask("clearData", task.ClearData)
	log.Fatal(exec.Run())
}

//xxl.Logger接口实现
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	fmt.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}

func (l *logger) Error(format string, a ...interface{}) {
	log.Println(fmt.Sprintf("自定义日志 - "+format, a...))
}
