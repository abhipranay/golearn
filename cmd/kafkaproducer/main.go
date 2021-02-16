package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/runtime/protoimpl"
	"math/rand"
	"os"
	"time"
)

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
			OldValue:      proto.Int32(1),
			NewValue:      proto.Int32(11),
		},
		Extras: nil,
	}
}

func dummyOfgStatusChangeMessage()  {
	
}

func main() {
	p, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": "localhost:9092",
		"client.id":         "oa_client",
		"acks":              "all"})

	if err != nil {
		fmt.Printf("Failed to create producer: %s\n", err)
		os.Exit(1)
	}
	r := rand.New(rand.NewSource(99))
	if err != nil {
		panic(err)
	}
	delivery_chan := make(chan kafka.Event, 10000)
	var value []byte
	topic := "order_status_change_test"
	for i := 0; i < 100; i++ {
		osc := dummyOrderStatusChange()
		if r.Int() % 2 == 0 {
			value, err = proto.Marshal(osc)
		} else {
			value, err = proto.Marshal(osc)
			//value, err = json.Marshal(osc)
		}
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          value},
			delivery_chan,
		)
		if err != nil {
			fmt.Println(err)
			continue
		}
		e := <-delivery_chan
		m := e.(*kafka.Message)
		o := &OrderStatusChange{}
		json.Unmarshal(m.Value, o)
		fmt.Println(o)
		time.Sleep(1 * time.Second)
	}
}
