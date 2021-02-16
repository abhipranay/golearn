package accounting

import (
	"fmt"
	"golearn/taskdispatcher/common"
	accounting "golearn/taskdispatcher/proto"
)

const (
	ORDER_EXT_UNPAID           = 1
	ORDER_EXT_RETURN_COMPLETED = 2
	ORDER_EXT_ESCROW_CREATED   = 3
	ORDER_EXT_CANCEL_COMPLETED = 4
)

// Rule tells conditions about order paid or created
func OrderPaidOrCreatedRule(m interface{}) bool {
	osc, ok := m.(*accounting.OrderStatusChange)
	if !ok {
		return false
	}
	newStatus := osc.GetStatusChange().GetNewValue()
	return newStatus == ORDER_EXT_UNPAID || newStatus == ORDER_EXT_ESCROW_CREATED
}

type SpexCallRequest struct {
	ShopId *int32
}

type SomeSpexCallTransformer struct{}

func (s SomeSpexCallTransformer) Transform(message common.MetaMessage) interface{} {
	orderStatusChangeMessage, _ := message.Value.(*accounting.OrderStatusChange)
	r := &SpexCallRequest{ShopId: orderStatusChangeMessage.ShopId}
	fmt.Println(r)
	return r
}

func GetAccountingEscrowCreatedRouter() common.Router {
	rb := &common.DefaultRouterBuilder{}
	kaleForwarder := &common.KaleForwarder{}
	ruleEngine := common.NewDispatcherRuleEngine()
	ruleEngine.AddRule(OrderPaidOrCreatedRule)
	kaleForwarder.Configure(&SomeSpexCallTransformer{})
	rb.AddForwarder(kaleForwarder)
	rb.SetSupportedMessageType("OrderStatusChange")
	rb.SetRuleEngine(ruleEngine)
	return rb.Build()
}
