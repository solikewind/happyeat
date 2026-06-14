package enum

import (
	"database/sql/driver"
	"fmt"
)

type SettlementStatus string

const (
	SettlementStatusUnsettled SettlementStatus = "UNSETTLED" // 未结账
	SettlementStatusSettled   SettlementStatus = "SETTLED"   // 已结账
)

func (s SettlementStatus) String() string {
	return string(s)
}

func (s SettlementStatus) Values() []string {
	return []string{
		string(SettlementStatusUnsettled),
		string(SettlementStatusSettled),
	}
}

func (s SettlementStatus) Value() (driver.Value, error) {
	return string(s), nil
}

func (s *SettlementStatus) Scan(value any) error {
	if value == nil {
		return nil
	}
	v, ok := value.(string)
	if !ok {
		if b, ok := value.([]byte); ok {
			v = string(b)
		} else {
			return fmt.Errorf("invalid type for SettlementStatus: %T", value)
		}
	}
	*s = SettlementStatus(v)
	return nil
}
