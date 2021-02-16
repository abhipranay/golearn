package ruleengine

type RuleEngine interface {
	AddRule(rule Rule)
	Evaluate(message interface{}) bool
}

type Rule func(message interface{}) bool

type engine struct {
	rules []Rule
}

func (e *engine) AddRule(rule Rule) {
	e.rules = append(e.rules, rule)
}

func (e *engine) Evaluate(message interface{}) bool {
	for _, rule := range e.rules {
		if !rule(message) {
			return false
		}
	}
	return true
}

func NewRuleEngine() RuleEngine {
	return &engine{rules: make([]Rule, 0)}
}