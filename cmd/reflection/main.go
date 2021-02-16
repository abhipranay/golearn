package main

import (
	"encoding/base64"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/cenkalti/backoff"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"reflect"
	"time"
)

const (
	ORDER_EXT_UNPAID           = 1
	ORDER_EXT_RETURN_COMPLETED = 2
	ORDER_EXT_ESCROW_CREATED   = 3
	ORDER_EXT_CANCEL_COMPLETED = 4
)

type Test1 struct {
	Name string `json:"name,omitempty"`
}

type Test2 struct {
	Email string `json:"email"`
	Age   int    `json:"age"`
}

func create(a interface{}, payload []byte) (string, interface{}, error) {
	var typ reflect.Type
	aType := reflect.TypeOf(a)
	aKind := aType.Kind()
	if aKind == reflect.Ptr {
		typ = aType.Elem()
	} else {
		typ = aType
	}

	f := reflect.New(typ).Interface()

	err := proto.Unmarshal(payload, f.(proto.Message))
	if err != nil {
		return "", nil, err
	}
	return reflect.TypeOf(f).String(), f, nil
}

func rawSarama(val []byte) *sarama.ConsumerMessage {
	return &sarama.ConsumerMessage{
		Headers:        nil,
		Timestamp:      time.Now(),
		BlockTimestamp: time.Now(),
		Key:            []byte("hello"),
		Value:          val,
		Topic:          "test1",
		Partition:      1,
		Offset:         100,
	}
}

func dummyOrderStatusChange() *OrderStatusChange {
	return &OrderStatusChange{
		state:         protoimpl.MessageState{},
		sizeCache:     0,
		unknownFields: nil,
		OrderId:       proto.Int64(1),
		UserId:        proto.Int32(2),
		ShopId:        proto.Int32(3),
		IsGroupBuy:    proto.Bool(false),
		StatusChange: &IntPair{
			state:         protoimpl.MessageState{},
			sizeCache:     0,
			unknownFields: nil,
			OldValue:      proto.Int32(ORDER_EXT_UNPAID),
			NewValue:      proto.Int32(ORDER_EXT_ESCROW_CREATED),
		},
		Extras: nil,
	}
}

func Base64Decode(message []byte) (b []byte, err error) {
	var l int
	b = make([]byte, base64.StdEncoding.DecodedLen(len(message)))
	l, err = base64.StdEncoding.Decode(b, message)
	if err != nil {
		return
	}
	return b[:l], nil
}

func RetryCheck() error {
	i := 1
	err := backoff.Retry(func() error {
		fmt.Println(i)
		i++
		if i == 9 {
			return nil
		}
		return fmt.Errorf("hello")
	}, backoff.WithMaxRetries(
		&backoff.ConstantBackOff{
			Interval: 10 * time.Second,
		}, uint64(10)),
	)
	return err
}

type A struct {
	val int
}

func unwrap(ai interface{})  {
	var m A
	m, ok := ai.(A)
	fmt.Printf("struct: %v\n", ok)
	if !ok {
		x, i := ai.(*A)
		if i {
			m = *x
		}
		fmt.Printf("ptr: %v\n", i)
	}
	fmt.Printf("m->val: %d\n", m.val)
}

func main() {
	a := A{5}
	var ai interface{}
	ai = a
	unwrap(ai)
}
