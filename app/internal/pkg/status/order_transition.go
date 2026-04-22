package status

import (
	"fmt"
	"strings"

	"github.com/solikewind/happyeat/common/consts/enum"
)

// 新建订单持久化时的初始状态请使用 enum.OrderStatusCreated（与状态机 NONE+TriggerCreate 结果一致，避免下单热路径再跑一遍内存状态机）。

// EntStatusToMachine 与 ent 枚举一致，直接作为状态机当前状态（大写）。
func EntStatusToMachine(s enum.OrderStatus) string {
	return strings.TrimSpace(string(s))
}

// MachineToEntStatus 状态机状态字符串转 ent 枚举。
func MachineToEntStatus(machine string) (enum.OrderStatus, error) {
	u := strings.ToUpper(strings.TrimSpace(machine))
	switch u {
	case OrderStatusCreated:
		return enum.OrderStatusCreated, nil
	case OrderStatusPaid:
		return enum.OrderStatusPaid, nil
	case OrderStatusPreparing:
		return enum.OrderStatusPreparing, nil
	case OrderStatusCompleted:
		return enum.OrderStatusCompleted, nil
	case OrderStatusCancelled:
		return enum.OrderStatusCancelled, nil
	default:
		return "", fmt.Errorf("未知状态机状态: %s", machine)
	}
}

// ParseAPIStatus 校验接口入参，统一为与 enum 一致的大写状态（兼容 created / CREATED）。
func ParseAPIStatus(s string) (string, error) {
	t := strings.TrimSpace(strings.ToLower(s))
	switch t {
	case "created":
		return OrderStatusCreated, nil
	case "paid":
		return OrderStatusPaid, nil
	case "preparing":
		return OrderStatusPreparing, nil
	case "completed":
		return OrderStatusCompleted, nil
	case "cancelled":
		return OrderStatusCancelled, nil
	default:
		return "", fmt.Errorf("无效的订单状态: %s", s)
	}
}

// ResolveTrigger 根据当前状态与目标状态返回 stateless 触发器。
func ResolveTrigger(from, to string) (string, error) {
	type pair struct{ f, t string }
	m := map[pair]string{
		{OrderStatusCreated, OrderStatusPaid}:        TriggerPay,
		{OrderStatusCreated, OrderStatusPreparing}:   TriggerPrepare,
		{OrderStatusCreated, OrderStatusCancelled}:    TriggerCancel,
		{OrderStatusPaid, OrderStatusPreparing}:      TriggerPrepare,
		{OrderStatusPaid, OrderStatusCancelled}:      TriggerCancel,
		{OrderStatusPreparing, OrderStatusCompleted}: TriggerComplete,
		{OrderStatusPreparing, OrderStatusCancelled}: TriggerCancel,
	}
	tr, ok := m[pair{from, to}]
	if !ok {
		return "", fmt.Errorf("当前状态 %q 不允许变更为 %q", from, to)
	}
	return tr, nil
}
