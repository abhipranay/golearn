package common

import (
	"github.com/Shopify/sarama"
	"time"
)

type Meta struct {
	MessageType string // OrderStatusChange
}

type MetaMessage struct {
	Meta      Meta // contains meta information like message type for message
	Topic     string
	Partition int64
	Offset    int64
	Timestamp time.Time
	Value     interface{} // contains concrete message struct like OrderStatusChange
}

type MessageTypeRegistry struct {
	register map[string]interface{} // contains mapping of topic -> concrete message struct
}

func (mr *MessageTypeRegistry) AddMessageType(topic string, message interface{}) {}

type KaleTaskConfig struct {
	Retry             int
	BackOffInterval   int
	BackOffMultiplier int
	Timeout           int
}

// Will be output of KaleMessageTransformers
type KaleTask struct {
	Command      string
	Request      interface{}
	RequestParam string
	KaleTaskConfig
}

type Rule func(message interface{}) bool

type KaleForwarder struct{}

func (k *KaleForwarder) Configure(transformer MessageTransformer, options ...interface{}) {
	panic("implement me")
}

func (k *KaleForwarder) Send(message MetaMessage) {
	panic("implement me")
}

type DispatcherRuleEngine struct {
	rules []Rule
}

func (d *DispatcherRuleEngine) AddRule(rule Rule) {
	if d.rules == nil {
		d.rules = make([]Rule, 0)
	}
	d.rules = append(d.rules, rule)
}

func (d *DispatcherRuleEngine) Evaluate(message interface{}) bool {
	for _, rule := range d.rules {
		if !rule(message) {
			return false
		}
	}
	return true
}

type JsonMessageTypeFinder struct {
	registry MessageTypeRegistry
}

func (r *JsonMessageTypeFinder) GetMessageType(rawMessage *sarama.ConsumerMessage) (string, interface{}, error) {
	panic("implement me")
}

type JsonMessageHandler struct {
	messageTypeFinder MessageTypeFinder
}

func (o *JsonMessageHandler) Handle(rawMessage *sarama.ConsumerMessage) error {
	panic("implement me")
}

type RouterImpl struct {
	ruleEngine RuleEngine
	forwarders []Forwarder
	msgType string
}

func (r RouterImpl) GetMessageType() string {
	panic("implement me")
}

func (r RouterImpl) Send(message MetaMessage) {
	panic("implement me")
}

type DefaultRouterBuilder struct {
	router RouterImpl
}

func (d *DefaultRouterBuilder) SetRuleEngine(ruleEngine RuleEngine) {
	d.router.ruleEngine = ruleEngine
}

func (d *DefaultRouterBuilder) AddForwarder(forwarder Forwarder) {
	d.router.forwarders = append(d.router.forwarders, forwarder)
}

func (d *DefaultRouterBuilder) SetSupportedMessageType(messageType string) {
	d.router.msgType = messageType
}

func (d *DefaultRouterBuilder) Build() Router {
	return d.router
}

