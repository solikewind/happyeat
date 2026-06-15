// Package order 提供订单与订单明细的 data 逻辑。
package order

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"entgo.io/ent/dialect/sql"
	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/dal/model/ent"
	entorder "github.com/solikewind/happyeat/dal/model/ent/order"
	entorderitem "github.com/solikewind/happyeat/dal/model/ent/orderitem"
	enttable "github.com/solikewind/happyeat/dal/model/ent/table"
)

// Order 订单数据访问。
type Order struct {
	c *ent.Client
}

// NewOrder 创建 Order。
func NewOrder(c *ent.Client) *Order {
	return &Order{c: c}
}

// ItemInput 订单项入参（下单时从菜单快照）
type ItemInput struct {
	MenuID    uint64
	MenuName  string
	Quantity  int
	UnitPrice int64
	SpecInfo  string
}

// CreateOrderInput 创建订单入参。
type CreateOrderInput struct {
	OrderType    enum.OrderType // dine_in | takeaway
	TableID      *uint64        // 堂食时必填，外带为 nil
	Items        []ItemInput    // 至少一项
	TotalAmount  int64
	ActualAmount int64 // 实收（分）；未收可为 0
	Remark       string
	Status       enum.OrderStatus // 初始状态，新建订单为 OrderStatusCreated（与状态机 NONE→create 一致）
}

// genOrderNo 生成订单号（简单示例：ORD+毫秒时间戳+3位随机）
func genOrderNo() string {
	return fmt.Sprintf("ORD%d%03d", time.Now().UnixMilli(), rand.Intn(1000))
}

// normalizeOrderStatus 统一为 ent 枚举大写（兼容历史/错误路径传入的 created 等小写）。
func normalizeOrderStatus(s enum.OrderStatus) enum.OrderStatus {
	u := strings.ToUpper(strings.TrimSpace(string(s)))
	switch u {
	case "CREATED":
		return enum.OrderStatusCreated
	case "PAID":
		return enum.OrderStatusPaid
	case "PREPARING":
		return enum.OrderStatusPreparing
	case "COMPLETED":
		return enum.OrderStatusCompleted
	case "CANCELLED":
		return enum.OrderStatusCancelled
	default:
		return s
	}
}

// Create 创建订单及明细（事务内）。业务规则（order_type、堂食必填桌台等）由调用方 Logic 校验。
func (o *Order) Create(ctx context.Context, in CreateOrderInput) (*ent.Order, error) {
	if len(in.Items) == 0 {
		return nil, fmt.Errorf("items required")
	}
	if in.Status == "" {
		return nil, fmt.Errorf("status required")
	}

	tx, err := o.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	create := tx.Order.Create().
		SetOrderNo(genOrderNo()).
		SetOrderType(in.OrderType).
		SetStatus(normalizeOrderStatus(in.Status)).
		SetTotalAmount(in.TotalAmount).
		SetActualAmount(in.ActualAmount)
	if in.TableID != nil && *in.TableID > 0 {
		create = create.SetTableID(*in.TableID)
	}
	if in.Remark != "" {
		create = create.SetRemark(in.Remark)
	}

	entOrder, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	for i, item := range in.Items {
		itemCreate := tx.OrderItem.Create().
			SetOrderID(entOrder.ID).
			SetMenuName(item.MenuName).
			SetQuantity(item.Quantity).
			SetUnitPrice(item.UnitPrice).
			SetAmount(item.UnitPrice * int64(item.Quantity)).
			SetSort(uint32(i))
		if item.MenuID > 0 {
			itemCreate = itemCreate.SetMenuID(item.MenuID)
		}
		if item.SpecInfo != "" {
			itemCreate = itemCreate.SetSpecInfo(item.SpecInfo)
		}
		if _, err = itemCreate.Save(ctx); err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return o.c.Order.Query().Where(entorder.IDEQ(entOrder.ID)).WithTable(func(q *ent.TableQuery) {
		q.WithCategory()
	}).WithItems().Only(ctx)
}

// GetByID 按 ID 获取订单（含 table、items）。
func (o *Order) GetByID(ctx context.Context, id uint64) (*ent.Order, error) {
	return o.c.Order.Query().
		Where(entorder.IDEQ(id)).
		WithTable(func(q *ent.TableQuery) {
			q.WithCategory()
		}).
		WithItems().
		Only(ctx)
}

// ListOrdersFilter 列表筛选。
type ListOrdersFilter struct {
	Status    enum.OrderStatus   // 单状态
	Statuses  []enum.OrderStatus // 多状态（与 Status 二选一，Statuses 优先）
	OrderType enum.OrderType
	TableID   *uint64
	OrderNo   string // 订单号模糊搜索
	Offset    int
	Limit     int
}

// List 分页列出订单（含 table、items），返回列表与总数。
func (o *Order) List(ctx context.Context, f ListOrdersFilter) ([]*ent.Order, int64, error) {
	q := o.c.Order.Query().WithTable(func(q *ent.TableQuery) {
		q.WithCategory()
	}).WithItems()
	if len(f.Statuses) > 0 {
		q = q.Where(entorder.StatusIn(f.Statuses...))
	} else if f.Status != "" {
		q = q.Where(entorder.StatusEQ(f.Status))
	}
	if f.OrderType != "" {
		q = q.Where(entorder.OrderTypeEQ(f.OrderType))
	}
	if f.TableID != nil && *f.TableID > 0 {
		q = q.Where(entorder.HasTableWith(enttable.IDEQ(*f.TableID)))
	}
	if f.OrderNo != "" {
		q = q.Where(entorder.OrderNoContainsFold(f.OrderNo))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if f.Limit <= 0 {
		f.Limit = 10
	}

	list, err := q.Order(entorder.ByID(sql.OrderDesc())).Limit(f.Limit).Offset(f.Offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return list, int64(total), nil
}

// UpdateStatus 更新订单状态；actualAmount 非 nil 时同时更新实收（分）。
func (o *Order) UpdateStatus(ctx context.Context, id uint64, status enum.OrderStatus, actualAmount *int64) error {
	update := o.c.Order.UpdateOneID(id).SetStatus(normalizeOrderStatus(status))
	if actualAmount != nil {
		if *actualAmount < 0 {
			return fmt.Errorf("actual_amount must be >= 0")
		}
		update = update.SetActualAmount(*actualAmount)
	}
	_, err := update.Save(ctx)
	return err
}

// AddItems 为订单追加菜单项，并重算 total_amount（actual_amount 仅在与 total_amount 相等时联动更新）。
func (o *Order) AddItems(ctx context.Context, id uint64, items []ItemInput) (*ent.Order, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items required")
	}
	tx, err := o.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	current, err := tx.Order.Query().Where(entorder.IDEQ(id)).WithItems().Only(ctx)
	if err != nil {
		return nil, err
	}
	existingItems, _ := current.Edges.ItemsOrErr()
	sortBase := len(existingItems)
	var appendAmount int64
	for i, item := range items {
		create := tx.OrderItem.Create().
			SetOrderID(id).
			SetMenuName(item.MenuName).
			SetQuantity(item.Quantity).
			SetUnitPrice(item.UnitPrice).
			SetAmount(item.UnitPrice * int64(item.Quantity)).
			SetSort(uint32(sortBase + i))
		if item.MenuID > 0 {
			create = create.SetMenuID(item.MenuID)
		}
		if item.SpecInfo != "" {
			create = create.SetSpecInfo(item.SpecInfo)
		}
		if _, err = create.Save(ctx); err != nil {
			return nil, err
		}
		appendAmount += item.UnitPrice * int64(item.Quantity)
	}

	newTotal := current.TotalAmount + appendAmount
	update := tx.Order.UpdateOneID(id).SetTotalAmount(newTotal)
	if current.ActualAmount == current.TotalAmount {
		update = update.SetActualAmount(current.ActualAmount + appendAmount)
	}
	if _, err = update.Save(ctx); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return o.GetByID(ctx, id)
}

// ReplaceItems 按提交列表整体替换订单明细，并重算 total_amount。
func (o *Order) ReplaceItems(ctx context.Context, id uint64, items []ItemInput) (*ent.Order, error) {
	if len(items) == 0 {
		return nil, fmt.Errorf("items required")
	}
	tx, err := o.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	current, err := tx.Order.Query().Where(entorder.IDEQ(id)).Only(ctx)
	if err != nil {
		return nil, err
	}

	// 先删旧明细，再按提交内容重建，保证“删除”在编辑后真实生效。
	if _, err = tx.OrderItem.Delete().Where(entorderitem.OrderIDEQ(id)).Exec(ctx); err != nil {
		return nil, err
	}

	var newTotal int64
	for i, item := range items {
		lineAmount := item.UnitPrice * int64(item.Quantity)
		create := tx.OrderItem.Create().
			SetOrderID(id).
			SetMenuName(item.MenuName).
			SetQuantity(item.Quantity).
			SetUnitPrice(item.UnitPrice).
			SetAmount(lineAmount).
			SetSort(uint32(i))
		if item.MenuID > 0 {
			create = create.SetMenuID(item.MenuID)
		}
		if item.SpecInfo != "" {
			create = create.SetSpecInfo(item.SpecInfo)
		}
		if _, err = create.Save(ctx); err != nil {
			return nil, err
		}
		newTotal += lineAmount
	}

	update := tx.Order.UpdateOneID(id).SetTotalAmount(newTotal)
	if current.ActualAmount == current.TotalAmount {
		update = update.SetActualAmount(newTotal)
	}
	if _, err = update.Save(ctx); err != nil {
		return nil, err
	}
	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return o.GetByID(ctx, id)
}
