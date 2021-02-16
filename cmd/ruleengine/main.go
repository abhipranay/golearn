package main

import (
	"fmt"
	"golearn/ruleengine"
)

type Poll1 struct {
	V1 int
	V2 int
}

type Poll struct {
	V1 int
	V2 int
}

func main() {
	re := ruleengine.NewRuleEngine()
	re.AddRule(CheckV1)

	res := re.Evaluate(&Poll{
		V1: 2,
		V2: 2,
	})
	fmt.Println(res)
}

func CheckV1(message interface{}) bool {
	//fmt.Println(reflect.TypeOf(message).ConvertibleTo(reflect.TypeOf(&Poll{})))
	m, ok := message.(*Poll)
	if !ok {
		return false
	}
	return m.V1 == 1
}
