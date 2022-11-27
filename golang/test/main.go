package main

import (
	"fmt"
	"sync"
	"time"
)

var sumTime time.Duration

var group sync.WaitGroup
var lock sync.Mutex

func main() {
	now := time.Now()
	for i := 0; i < 200; i++ {
		group.Add(1)
		go writeFile()
	}
	group.Wait()
	fmt.Printf("合计耗时：%v\n", sumTime)
	fmt.Printf("程序耗时:%v\n", time.Since(now))
}
func writeFile() {
	defer group.Done()
	for i := 0; i < 100; i++ {
		lock.Lock()
		sumTime += 1 * time.Second
		lock.Unlock()
	}

}
