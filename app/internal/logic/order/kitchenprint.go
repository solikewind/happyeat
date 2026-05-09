package order

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/zeromicro/go-zero/core/logx"
)

// SyncPrintKitchen 同步调用商鹏厨房单；用于「手动打印」接口，失败将返回给调用方。
func SyncPrintKitchen(ctx context.Context, svcCtx *svc.ServiceContext, e *ent.Order, banner string) error {
	if svcCtx == nil || e == nil {
		return errors.New("内部错误")
	}
	if svcCtx.Spyun == nil {
		return errors.New("商鹏云打印未启用或未配置")
	}
	content := formatKitchenTicket(e, banner)
	_, err := svcCtx.Spyun.PrintOrder(ctx, "", content, 1)
	return err
}

// scheduleKitchenPrint 异步提交商鹏打印；nil 客户端或未启用时不做事。失败只记日志，不影响订单接口结果。
func scheduleKitchenPrint(svcCtx *svc.ServiceContext, e *ent.Order, banner string) {
	if svcCtx == nil || svcCtx.Spyun == nil || e == nil {
		return
	}
	content := formatKitchenTicket(e, banner)
	orderNo := e.OrderNo
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
		defer cancel()
		reply, err := svcCtx.Spyun.PrintOrder(ctx, "", content, 1)
		if err != nil {
			logx.WithContext(ctx).Errorf("商鹏厨房单打印失败 order_no=%s: %v", orderNo, err)
			return
		}
		if reply != nil && reply.ID != "" {
			logx.WithContext(ctx).Infof("商鹏厨房单已受理 order_no=%s print_id=%s", orderNo, reply.ID)
		}
	}()
}

func formatKitchenTicket(e *ent.Order, banner string) string {
	var b strings.Builder
	banner = strings.TrimSpace(banner)
	if banner != "" {
		b.WriteString("<C><B>")
		b.WriteString(escapeSpyunText(banner))
		b.WriteString("</B></C><BR>")
	}
	b.WriteString("订单: ")
	b.WriteString(escapeSpyunText(e.OrderNo))
	b.WriteString("<BR>")
	b.WriteString("下单: ")
	b.WriteString(e.CreatedAt.Format("2006-01-02 15:04"))
	b.WriteString("<BR>")
	switch string(e.OrderType) {
	case "dine_in":
		b.WriteString("类型: 堂食<BR>")
		tbl, err := e.Edges.TableOrErr()
		if err == nil && tbl != nil {
			b.WriteString("桌台: ")
			b.WriteString(escapeSpyunText(tbl.Code))
			cat, _ := tbl.Edges.CategoryOrErr()
			if cat != nil && strings.TrimSpace(cat.Name) != "" {
				b.WriteString(" (")
				b.WriteString(escapeSpyunText(strings.TrimSpace(cat.Name)))
				b.WriteString(")")
			}
			b.WriteString("<BR>")
		}
	case "takeaway":
		b.WriteString("类型: 外带<BR>")
	default:
		b.WriteString("类型: ")
		b.WriteString(escapeSpyunText(string(e.OrderType)))
		b.WriteString("<BR>")
	}
	b.WriteString("--------------------------------<BR>")

	items, _ := e.Edges.ItemsOrErr()
	for _, it := range items {
		if it == nil {
			continue
		}
		line := fmt.Sprintf("%s x%d  %s元",
			escapeSpyunText(it.MenuName),
			it.Quantity,
			fmtYuan(it.UnitPrice),
		)
		b.WriteString(line)
		b.WriteString("<BR>")
		if it.SpecInfo != nil && strings.TrimSpace(*it.SpecInfo) != "" {
			b.WriteString("  规格: ")
			b.WriteString(escapeSpyunText(strings.TrimSpace(*it.SpecInfo)))
			b.WriteString("<BR>")
		}
	}
	b.WriteString("--------------------------------<BR>")
	b.WriteString(fmt.Sprintf("合计: %s元<BR>", fmtYuan(e.TotalAmount)))
	b.WriteString(fmt.Sprintf("实收: %s元<BR>", fmtYuan(e.ActualAmount)))
	if e.Remark != nil && strings.TrimSpace(*e.Remark) != "" {
		b.WriteString("备注: ")
		b.WriteString(escapeSpyunText(strings.TrimSpace(*e.Remark)))
		b.WriteString("<BR>")
	}
	b.WriteString("<CUT>")
	return b.String()
}

func fmtYuan(cents int64) string {
	return fmt.Sprintf("%.2f", float64(cents)/100)
}

// 商鹏内容为文本指令混排，避免顾客备注等破坏标签：去掉尖括号与换行。
func escapeSpyunText(s string) string {
	s = strings.ReplaceAll(s, "<", "＜")
	s = strings.ReplaceAll(s, ">", "＞")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}
