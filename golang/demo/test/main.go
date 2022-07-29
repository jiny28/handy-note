package main

import "fmt"

type TagValue struct {
	name  string
	value interface{}
}

func main() {

	b := int64(1)

	a := TagValue{name: "a", value: b}
	fmt.Println("输出", a.value)

}
