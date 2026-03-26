package enum

import (
	"database/sql/driver"
	"fmt"
)

type OrderStatus string

const (
	OrderStatusCreated   OrderStatus = "CREATED"
	OrderStatusPaid      OrderStatus = "PAID"
	OrderStatusPreparing OrderStatus = "PREPARING"
	OrderStatusCompleted OrderStatus = "COMPLETED"
	OrderStatusCancelled OrderStatus = "CANCELLED"
)

func (s OrderStatus) String() string {
	return string(s)
}

func (s OrderStatus) Values() []string {
	return []string{
		string(OrderStatusCreated),
		string(OrderStatusPaid),
		string(OrderStatusPreparing),
		string(OrderStatusCompleted),
		string(OrderStatusCancelled),
	}
}

func (s OrderStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *OrderStatus) Scan(value any) error {
	if value == nil {
		return nil
	}
	v, ok := value.(string)
	if !ok {
		// 部分驱动可能返回 []byte
		if b, ok := value.([]byte); ok {
			v = string(b)
		} else {
			return fmt.Errorf("invalid type for OrderStatus: %T", value)
		}
	}
	*s = OrderStatus(v)
	return nil
}
