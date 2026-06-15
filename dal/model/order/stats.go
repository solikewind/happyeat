package order

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/solikewind/happyeat/common/consts/enum"
	"github.com/solikewind/happyeat/dal/model/ent"
	entorder "github.com/solikewind/happyeat/dal/model/ent/order"
	entschema "github.com/solikewind/happyeat/dal/model/ent/schema"
)

// StatsDateRange 统计日期闭区间 [Start, End]（CST 自然日）。
type StatsDateRange struct {
	Start        time.Time
	End          time.Time
	EndExclusive time.Time
}

// DailyStatsPoint 按日汇总。
type DailyStatsPoint struct {
	Date          time.Time
	OrderCount    int
	Revenue       int64
	Receivable    int64
	ActualRevenue int64
	ItemCount     int
	DineInCount   int
	TakeawayCount int
}

// StatsSummary 区间汇总。
type StatsSummary struct {
	OrderCount    int
	Revenue       int64
	Receivable    int64
	ActualRevenue int64
	ItemCount     int
	DineInCount   int
	TakeawayCount int
}

// MenuSalesRow 菜品销量行。
type MenuSalesRow struct {
	MenuName string
	SpecInfo string
	Quantity int
	Amount   int64
}

const statsDateLayout = "2006-01-02"

// ParseStatsDateRange 解析 YYYY-MM-DD；缺省均为今天；仅 start 时 end=start。
func ParseStatsDateRange(startStr, endStr string) (StatsDateRange, error) {
	today := dateOnlyCST(time.Now())
	start := today
	end := today

	if startStr != "" {
		t, err := parseStatsDate(startStr)
		if err != nil {
			return StatsDateRange{}, fmt.Errorf("start_date: %w", err)
		}
		start = t
	}
	if endStr != "" {
		t, err := parseStatsDate(endStr)
		if err != nil {
			return StatsDateRange{}, fmt.Errorf("end_date: %w", err)
		}
		end = t
	} else if startStr != "" {
		end = start
	}
	if end.Before(start) {
		start, end = end, start
	}
	return StatsDateRange{
		Start:        start,
		End:          end,
		EndExclusive: end.AddDate(0, 0, 1),
	}, nil
}

func parseStatsDate(s string) (time.Time, error) {
	t, err := time.ParseInLocation(statsDateLayout, s, entschema.CST)
	if err != nil {
		return time.Time{}, err
	}
	return dateOnlyCST(t), nil
}

func dateOnlyCST(t time.Time) time.Time {
	cst := t.In(entschema.CST)
	return time.Date(cst.Year(), cst.Month(), cst.Day(), 0, 0, 0, 0, entschema.CST)
}

func saleStatuses() []enum.OrderStatus {
	return []enum.OrderStatus{
		enum.OrderStatusPaid,
		enum.OrderStatusPreparing,
		enum.OrderStatusCompleted,
	}
}

func orderRevenue(o *ent.Order) int64 {
	if o.ActualAmount > 0 {
		return o.ActualAmount
	}
	return o.TotalAmount
}

func itemLineAmount(it *ent.OrderItem) int64 {
	if it.Amount > 0 {
		return it.Amount
	}
	return it.UnitPrice * int64(it.Quantity)
}

func (o *Order) listSaleOrdersInRange(ctx context.Context, r StatsDateRange) ([]*ent.Order, error) {
	return o.c.Order.Query().
		Where(
			entorder.StatusIn(saleStatuses()...),
			entorder.CreatedAtGTE(r.Start),
			entorder.CreatedAtLT(r.EndExclusive),
		).
		WithItems().
		All(ctx)
}

// AggregateSummary 区间订单汇总。
func (o *Order) AggregateSummary(ctx context.Context, r StatsDateRange) (StatsSummary, error) {
	orders, err := o.listSaleOrdersInRange(ctx, r)
	if err != nil {
		return StatsSummary{}, err
	}
	var sum StatsSummary
	for _, ord := range orders {
		sum.OrderCount++
		rev := orderRevenue(ord)
		sum.Revenue += rev
		sum.Receivable += ord.TotalAmount
		sum.ActualRevenue += rev
		if ord.OrderType == enum.OrderTypeDineIn {
			sum.DineInCount++
		} else {
			sum.TakeawayCount++
		}
		items, _ := ord.Edges.ItemsOrErr()
		for _, it := range items {
			sum.ItemCount += it.Quantity
		}
	}
	return sum, nil
}

// AggregateDaily 按 CST 自然日分组汇总，区间内无数据的日期补零。
func (o *Order) AggregateDaily(ctx context.Context, r StatsDateRange) ([]DailyStatsPoint, error) {
	orders, err := o.listSaleOrdersInRange(ctx, r)
	if err != nil {
		return nil, err
	}
	byDay := make(map[time.Time]*DailyStatsPoint)
	for _, ord := range orders {
		day := dateOnlyCST(ord.CreatedAt)
		p, ok := byDay[day]
		if !ok {
			p = &DailyStatsPoint{Date: day}
			byDay[day] = p
		}
		p.OrderCount++
		rev := orderRevenue(ord)
		p.Revenue += rev
		p.Receivable += ord.TotalAmount
		p.ActualRevenue += rev
		if ord.OrderType == enum.OrderTypeDineIn {
			p.DineInCount++
		} else {
			p.TakeawayCount++
		}
		items, _ := ord.Edges.ItemsOrErr()
		for _, it := range items {
			p.ItemCount += it.Quantity
		}
	}
	out := make([]DailyStatsPoint, 0)
	for d := r.Start; !d.After(r.End); d = d.AddDate(0, 0, 1) {
		if p, ok := byDay[d]; ok {
			out = append(out, *p)
			continue
		}
		out = append(out, DailyStatsPoint{Date: d})
	}
	return out, nil
}

// AggregateMenuSales 区间内按菜品+规格汇总销量。
func (o *Order) AggregateMenuSales(ctx context.Context, r StatsDateRange) ([]MenuSalesRow, error) {
	orders, err := o.listSaleOrdersInRange(ctx, r)
	if err != nil {
		return nil, err
	}
	type key struct {
		name string
		spec string
	}
	agg := make(map[key]*MenuSalesRow)
	for _, ord := range orders {
		items, _ := ord.Edges.ItemsOrErr()
		for _, it := range items {
			spec := ""
			if it.SpecInfo != nil {
				spec = strings.TrimSpace(*it.SpecInfo)
			}
			k := key{name: it.MenuName, spec: spec}
			row, ok := agg[k]
			if !ok {
				row = &MenuSalesRow{MenuName: it.MenuName, SpecInfo: spec}
				agg[k] = row
			}
			row.Quantity += it.Quantity
			row.Amount += itemLineAmount(it)
		}
	}
	out := make([]MenuSalesRow, 0, len(agg))
	for _, row := range agg {
		out = append(out, *row)
	}
	return out, nil
}
