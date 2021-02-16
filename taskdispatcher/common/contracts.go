package common

import (
	"github.com/Shopify/sarama"
)

type MessageTypeFinder interface {
	// (*sarama.ConsumerMessage) -> (messageType string, concrete message struct, error if not able to get concrete type)
	GetMessageType(rawMessage *sarama.ConsumerMessage) (string, interface{}, error)
}

// entry point for dispatcher. Receives a raw sarama.ConsumerMessage and figures out
// concrete type of message by consulting message registry and pass MetaMessage (which contains concrete type)
// to all the Router which subscribe for this message type
type MessageHandler interface {
	Handle(rawMessage *sarama.ConsumerMessage) error
}

type MessageHandlerBuilder interface {
	AddRouter(router Router)
	SetMessageTypeFinder(mtf MessageTypeFinder)
	Build() MessageHandler
}

// RuleEngine is a collection of rules and knows how to evaluate rules in it
type RuleEngine interface {
	AddRule(rule Rule)
	Evaluate(message interface{}) bool
}

// Router consists of a RuleEngine and Forwarder(s)
// It receives a MetaMessage -> Evaluates via RuleEngine -> True Result -> Pass to all Forwarder(s)
type Router interface {
	GetMessageType() string
	Send(message MetaMessage)
}

type RouterBuilder interface {
	SetRuleEngine(ruleEngine RuleEngine)
	AddForwarder(forwarder Forwarder)
	SetSupportedMessageType(messageType string)
	Build() Router
}

// This component helps build request which a Forwarder will send outside. Eg: Kale expects a
// Spex command request. So an implementation will help prepare Kale command request object
type MessageTransformer interface {
	Transform(message MetaMessage) interface{}
}

type Forwarder interface {
	Configure(options ...interface{})
	Send(message MetaMessage)
}

