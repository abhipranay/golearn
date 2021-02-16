package context

import (
	"context"
	"fmt"
	"time"
)

func lazy(t uint) bool {
	d := time.Duration(t) * time.Second
	time.Sleep(d)
	fmt.Println("lazy is done")
	return true
}

func Run(timeout, lt uint) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeout)*time.Second)
	callDone := make(chan bool)
	defer cancel()
	go func() {
		callDone <- lazy(lt)
	}()
	select {
	case <-ctx.Done():
		fmt.Println("done", time.Now())
		return ctx.Err()
	case <-callDone:
		fmt.Println("lazy", time.Now())
	}
	return nil
}

func Select() {
	chan1 := make(chan bool)
	chan2 := make(chan bool)

	go func(ch chan bool) {
		time.Sleep(4 * time.Second)
		ch <- true
	}(chan1)

	go func(ch chan bool) {
		time.Sleep(2 * time.Second)
		ch <- true
	}(chan2)

	select {
	case <-chan1:
		fmt.Println("Read from chan1")
	case <-chan2:
		fmt.Println("Read from chan2")
	}
}
