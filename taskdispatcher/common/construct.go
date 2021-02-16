package common

func NewDispatcherRuleEngine() *DispatcherRuleEngine {
	return &DispatcherRuleEngine{
		rules: make([]Rule, 0),
	}
}


