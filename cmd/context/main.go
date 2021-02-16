package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"os/signal"
	"runtime"
	"time"
)

func hello(ctx context.Context, msg string) {
	PrintMemUsage()
	select {
	case <-time.After(2 * time.Second):
	}
	if ctx.Err() != nil {
		return
	}
	fmt.Println(msg)
}

func createNewCtx(parent context.Context) {
	path := "/Users/abhipranay.chauhan/Downloads/test.jpg"
	data, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	ctx, cancel := context.WithCancel(context.WithValue(parent, "data", data))
	go func() {
		hello(ctx, "hello")
		cancel()
		PrintMemUsage()
	}()
}

func main() {
	run()
}

func run() {
	PrintMemUsage()
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Kill, os.Interrupt)
	ctx := context.Background()
	for i := 0; i < 5; i++ {
		go createNewCtx(ctx)
	}
	<-sig
	PrintMemUsage()
}

// PrintMemUsage outputs the current, total and OS memory being used. As well as the number
// of garage collection cycles completed.
func PrintMemUsage() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}