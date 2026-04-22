package status

// 订单状态：与 common/consts/enum、ent 存库一致（大写），stateless Configure 直接使用这些字符串。
const (
	OrderStatusNone      = "NONE"
	OrderStatusCreated   = "CREATED"
	OrderStatusPaid      = "PAID"
	OrderStatusPreparing = "PREPARING"
	OrderStatusCompleted = "COMPLETED"
	OrderStatusCancelled = "CANCELLED"
)

// 状态触发器（事件），小写动词即可，与状态字符串解耦。
const (
	TriggerCreate   = "create"
	TriggerPay      = "pay"
	TriggerPrepare  = "prepare"
	TriggerComplete = "complete"
	TriggerCancel   = "cancel"
)
