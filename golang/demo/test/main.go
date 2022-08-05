package main

import (
	"fmt"
	"runtime"
	"time"
)

type Fn func() error

type MyTicker struct {
	MyTick *time.Ticker
	Runner Fn
}

func NewMyTick(interval int, f Fn) *MyTicker {
	return &MyTicker{
		MyTick: time.NewTicker(time.Duration(interval) * time.Second),
		Runner: f,
	}
}

func (t *MyTicker) Start() {
	for {
		select {
		case <-t.MyTick.C:
			t.Runner()
		}
	}
}

var testP Fn = func() error {
	fmt.Println(" 滴答 1 次")
	return nil
}

func testPrint() {

}

func main() {
	fmt.Println(runtime.NumCPU())
	t := NewMyTick(1, testP)
	go t.Start()
	fmt.Println("start .")
	for {
		fmt.Println("for something.")
		time.Sleep(time.Second * 10)

	}
}
