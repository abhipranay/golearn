package main

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"time"
)

func genMetrics(c *DefaultCollector, done chan int)  {
	for true {
		c.messageCount.WithLabelValues("testing").Inc()
		fmt.Println("sent metrics")
		time.Sleep(2 * time.Second)
	}
	done <- 1
}

func main() {
	registry := prometheus.NewRegistry()
	done := make(chan int)
	mux := http.NewServeMux()
	srv := &http.Server{
		Handler: mux,
		Addr:    fmt.Sprintf(":%v", "12345"),
	}
	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			panic(err)
		}
	}()
	c := NewMetricsCollector(registry)
	go genMetrics(c, done)
	<-done
}