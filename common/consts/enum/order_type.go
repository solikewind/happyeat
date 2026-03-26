package enum

import (
	"database/sql/driver"
	"fmt"
)

type OrderType string

const (
	OrderTypeDineIn   OrderType = "dine_in"
	OrderTypeTakeaway OrderType = "takeaway"
)

func (o OrderType) String() string {
	return string(o)
}

func (o OrderType) Values() []string {
	return []string{
		string(OrderTypeDineIn),
		string(OrderTypeTakeaway),
	}
}

func (o OrderType) Value() (driver.Value, error) {
	return string(o), nil
}

func (o *OrderType) Scan(value any) error {
	if value == nil {
		return nil
	}
	v, ok := value.(string)
	if !ok {
		// 部分驱动可能返回 []byte
		if b, ok := value.([]byte); ok {
			v = string(b)
		} else {
			return fmt.Errorf("invalid type for OrderType: %T", value)
		}
	}
	*o = OrderType(v)
	return nil
}
