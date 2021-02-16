package main

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	messageResponseTimeHistogram = "message_response_time_histogram"
	messageDelayTimeHistogram    = "message_delay_time_histogram"
	messageCount                 = "message_count"
	messageErrorCount            = "message_error_count"
	transporterErrorCount        = "transporter_error_count"
	messageTypeTag               = "message_type"
	transporterTypeTag           = "transporter_type"
	errorCodeTag                 = "error_code"
)

var (
	delayMillisecondsBuckets    = []float64{100, 1000, 10000, 100000, 1000000, 10000000}
	responseMillisecondsBuckets = []float64{10, 100, 500, 1000, 2000, 4000}
)

type DefaultCollector struct {
	messageResponseTimeHistogram *prometheus.HistogramVec
	messageDelayTimeHistogram    *prometheus.HistogramVec
	messageCount                 *prometheus.CounterVec
	handleErrorCount             *prometheus.CounterVec
	transporterErrorCount        *prometheus.CounterVec
}

func register(registerer prometheus.Registerer) *DefaultCollector {
	counter := prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: messageCount,
		Help: "count of messages per event type",
	}, []string{messageTypeTag})
	registerer.MustRegister(counter)
	return &DefaultCollector{
		messageCount: counter,
	}
}

func (c DefaultCollector) RecordMessageResponseTime(messageType string, responseTime float64) {
	c.messageResponseTimeHistogram.WithLabelValues(messageType).Observe(responseTime)
}

func (c DefaultCollector) RecordMessageDelayTime(messageType string, delayTime float64) {
	o, err := c.messageDelayTimeHistogram.GetMetricWithLabelValues(messageType)
	if err != nil {
		return
	}
	o.Observe(delayTime)
}

func (c DefaultCollector) RecordMessageCount(messageType string) {
	c.messageCount.WithLabelValues(messageType).Inc()
}

func (c DefaultCollector) RecordErrorCount(messageType, errorCode string) {
	c.handleErrorCount.WithLabelValues(messageType, errorCode).Inc()
}

func (c DefaultCollector) RecordTransporterErrorCount(transporterType, errorCode string) {
	c.transporterErrorCount.WithLabelValues(transporterType, errorCode).Inc()
}

func NewMetricsCollector(registerer prometheus.Registerer) *DefaultCollector {
	return register(registerer)
}

