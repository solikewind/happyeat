package status

// 订单状态
const (
	OrderStatusNone      = "none"      // 无状态，初始状态
	OrderStatusCreated   = "created"   // 订单已创建，等待支付
	OrderStatusPaid      = "paid"      // 订单已支付，等待商家接单
	OrderStatusPreparing = "preparing" // 商家已接单，正在准备餐品
	OrderStatusCompleted = "completed" // 订单已完成
	OrderStatusCancelled = "cancelled" // 订单已取消
)

// 状态触发器（事件）
const (
	TriggerCreate   = "create"   // 创建订单
	TriggerPay      = "pay"      // 支付订单
	TriggerPrepare  = "prepare"  // 商家接单
	TriggerComplete = "complete" // 完成订单
	TriggerCancel   = "cancel"   // 取消订单
)
