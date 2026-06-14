package settlement

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/dal/model/ent"
	entorder "github.com/solikewind/happyeat/dal/model/ent/order"
	entsettlement "github.com/solikewind/happyeat/dal/model/ent/settlement"
)

// Settlement 结账单 data 层。
type Settlement struct {
	c *ent.Client
}

func NewSettlement(c *ent.Client) *Settlement {
	return &Settlement{c: c}
}

type ListFilter struct {
	Status       enum.SettlementStatus
	CustomerName string
	Offset       int
	Limit        int
}

func (s *Settlement) Create(ctx context.Context, customerName, remark string) (*ent.Settlement, error) {
	name := trimCustomerName(customerName)
	if name == "" {
		return nil, fmt.Errorf("customer_name required")
	}
	exists, err := s.c.Settlement.Query().
		Where(
			entsettlement.CustomerNameEQ(name),
			entsettlement.StatusEQ(enum.SettlementStatusUnsettled),
		).
		Exist(ctx)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, fmt.Errorf("该客户已有未结账的结账单")
	}

	create := s.c.Settlement.Create().
		SetCustomerName(name).
		SetStatus(enum.SettlementStatusUnsettled)
	if remark != "" {
		create = create.SetRemark(remark)
	}
	return create.Save(ctx)
}

func (s *Settlement) GetByID(ctx context.Context, id uint64) (*ent.Settlement, error) {
	return s.c.Settlement.Query().
		Where(entsettlement.IDEQ(id)).
		WithOrders(func(q *ent.OrderQuery) {
			q.WithTable(func(tq *ent.TableQuery) {
				tq.WithCategory()
			}).WithItems()
		}).
		Only(ctx)
}

func (s *Settlement) List(ctx context.Context, f ListFilter) ([]*ent.Settlement, int64, error) {
	q := s.c.Settlement.Query()
	if f.Status != "" {
		q = q.Where(entsettlement.StatusEQ(f.Status))
	}
	if f.CustomerName != "" {
		q = q.Where(entsettlement.CustomerNameContains(trimCustomerName(f.CustomerName)))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if f.Limit <= 0 {
		f.Limit = 10
	}

	list, err := q.
		Order(ent.Desc(entsettlement.FieldID)).
		Limit(f.Limit).
		Offset(f.Offset).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}
	return list, int64(total), nil
}

func (s *Settlement) AddOrder(ctx context.Context, settlementID, orderID uint64) (*ent.Settlement, error) {
	tx, err := s.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	st, err := tx.Settlement.Query().Where(entsettlement.IDEQ(settlementID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("结账单不存在")
		}
		return nil, err
	}
	if st.Status != enum.SettlementStatusUnsettled {
		return nil, fmt.Errorf("仅未结账的结账单可添加订单")
	}

	ord, err := tx.Order.Query().Where(entorder.IDEQ(orderID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, err
	}
	if ord.Status == enum.OrderStatusCancelled {
		return nil, fmt.Errorf("已取消的订单不能加入结账单")
	}
	if ord.SettlementID != nil && *ord.SettlementID > 0 {
		if *ord.SettlementID == settlementID {
			if err := tx.Commit(); err != nil {
				return nil, err
			}
			return s.GetByID(ctx, settlementID)
		}
		other, err := tx.Settlement.Get(ctx, *ord.SettlementID)
		if err != nil {
			return nil, err
		}
		if other.Status == enum.SettlementStatusUnsettled {
			return nil, fmt.Errorf("订单已在其他未结账结账单中")
		}
		return nil, fmt.Errorf("订单已在已结账结账单中，不能重复加入")
	}

	if _, err = tx.Order.UpdateOneID(orderID).SetSettlementID(settlementID).Save(ctx); err != nil {
		return nil, err
	}
	if err = recalcTotal(ctx, tx, settlementID); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, settlementID)
}

func (s *Settlement) RemoveOrder(ctx context.Context, settlementID, orderID uint64) (*ent.Settlement, error) {
	tx, err := s.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	st, err := tx.Settlement.Query().Where(entsettlement.IDEQ(settlementID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("结账单不存在")
		}
		return nil, err
	}
	if st.Status != enum.SettlementStatusUnsettled {
		return nil, fmt.Errorf("仅未结账的结账单可移除订单")
	}

	ord, err := tx.Order.Query().Where(entorder.IDEQ(orderID)).Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("订单不存在")
		}
		return nil, err
	}
	if ord.SettlementID == nil || *ord.SettlementID != settlementID {
		return nil, fmt.Errorf("订单不在该结账单中")
	}

	if _, err = tx.Order.UpdateOneID(orderID).ClearSettlementID().Save(ctx); err != nil {
		return nil, err
	}
	if err = recalcTotal(ctx, tx, settlementID); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, settlementID)
}

type SettleInput struct {
	ActualAmount int64
	Remark       string
}

func (s *Settlement) Settle(ctx context.Context, settlementID uint64, in SettleInput) (*ent.Settlement, error) {
	if in.ActualAmount < 0 {
		return nil, fmt.Errorf("实收金额不能为负")
	}

	tx, err := s.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	st, err := tx.Settlement.Query().
		Where(entsettlement.IDEQ(settlementID)).
		WithOrders().
		Only(ctx)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, fmt.Errorf("结账单不存在")
		}
		return nil, err
	}
	if st.Status != enum.SettlementStatusUnsettled {
		return nil, fmt.Errorf("结账单已结账")
	}

	orders, _ := st.Edges.OrdersOrErr()
	if len(orders) == 0 {
		return nil, fmt.Errorf("结账单中没有订单，无法结账")
	}

	if err = recalcTotal(ctx, tx, settlementID); err != nil {
		return nil, err
	}
	st, err = tx.Settlement.Get(ctx, settlementID)
	if err != nil {
		return nil, err
	}

	upd := tx.Settlement.UpdateOneID(settlementID).
		SetStatus(enum.SettlementStatusSettled).
		SetActualAmount(in.ActualAmount).
		SetSettledAt(time.Now())
	if in.Remark != "" {
		upd = upd.SetRemark(in.Remark)
	} else if st.Remark != nil {
		upd = upd.SetNillableRemark(st.Remark)
	}
	if _, err = upd.Save(ctx); err != nil {
		return nil, err
	}

	for _, ord := range orders {
		if ord.Status == enum.OrderStatusCancelled {
			continue
		}
		if _, err = tx.Order.UpdateOneID(ord.ID).
			SetActualAmount(ord.TotalAmount).
			Save(ctx); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, settlementID)
}

func recalcTotal(ctx context.Context, tx *ent.Tx, settlementID uint64) error {
	orders, err := tx.Order.Query().
		Where(
			entorder.SettlementIDEQ(settlementID),
			entorder.StatusNEQ(enum.OrderStatusCancelled),
		).
		All(ctx)
	if err != nil {
		return err
	}
	var total int64
	for _, o := range orders {
		total += o.TotalAmount
	}
	_, err = tx.Settlement.UpdateOneID(settlementID).SetTotalAmount(total).Save(ctx)
	return err
}

func (s *Settlement) CountOrders(ctx context.Context, settlementID uint64) (int, error) {
	return s.c.Order.Query().
		Where(
			entorder.SettlementIDEQ(settlementID),
			entorder.StatusNEQ(enum.OrderStatusCancelled),
		).
		Count(ctx)
}

// RecalcTotalForOrder 若订单挂在未结账结账单上，重算该结账单应收合计。
func (s *Settlement) RecalcTotalForOrder(ctx context.Context, orderID uint64) error {
	ord, err := s.c.Order.Get(ctx, orderID)
	if err != nil {
		return err
	}
	if ord.SettlementID == nil || *ord.SettlementID == 0 {
		return nil
	}
	st, err := s.c.Settlement.Get(ctx, *ord.SettlementID)
	if err != nil {
		return err
	}
	if st.Status != enum.SettlementStatusUnsettled {
		return nil
	}
	tx, err := s.c.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err = recalcTotal(ctx, tx, st.ID); err != nil {
		return err
	}
	return tx.Commit()
}

func trimCustomerName(name string) string {
	return strings.TrimSpace(name)
}
