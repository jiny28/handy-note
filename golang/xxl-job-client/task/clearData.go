package task

import (
	"context"
	"fmt"
	xxl "github.com/xxl-job/xxl-job-executor-go"
	"strconv"
)

func ClearData(cxt context.Context, param *xxl.RunReq) (msg string) {

	fmt.Println("clear data exec param " + param.ExecutorParams + ",index:" + strconv.FormatInt(param.BroadcastIndex, 10) + ",total:" + strconv.FormatInt(param.BroadcastTotal, 10))

	return
}
