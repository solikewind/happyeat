package status

import (
	"log"

	"github.com/qmuntal/stateless"
	"github.com/solikewind/happyeat/dal/model/ent"
)

// OrderStateMachine 封装 stateless 状态机
type OrderStateMachine struct {
	sm *stateless.StateMachine
}

func NewOrderStateMachine(currentStatus string, order ent.Order) *OrderStateMachine {
	sm := stateless.NewStateMachine(currentStatus) // 创建状态机

	// 创建 - > 支付 -> 制作 -> 完成
	// 创建 - > （无需支付） -> 制作 -> 完成
	// 定义状态转换
	// 订单状态转换
	// 无状态 -> 已创建
	sm.Configure(OrderStatusNone).
		Permit(TriggerCreate, OrderStatusCreated).
		Permit(TriggerCancel, OrderStatusCancelled)

	// 已创建 -> 已支付 -> 商家接单 -> 完成
	sm.Configure(OrderStatusCreated).
		Permit(TriggerPay, OrderStatusPaid).          // 需要支付
		Permit(TriggerPrepare, OrderStatusPreparing). // 无需支付，直接接单
		Permit(TriggerCancel, OrderStatusCancelled)

	// 已支付 -> 准备中
	sm.Configure(OrderStatusPaid).
		Permit(TriggerPrepare, OrderStatusPreparing).
		Permit(TriggerCancel, OrderStatusCancelled)

	// 已接单 -> 完成
	sm.Configure(OrderStatusPreparing).
		Permit(TriggerComplete, OrderStatusCompleted).
		Permit(TriggerCancel, OrderStatusCancelled)

	graph := sm.ToGraph()
	log.Printf("Order State Machine Graph: %v", graph)

	return &OrderStateMachine{
		sm: sm,
	}
}
