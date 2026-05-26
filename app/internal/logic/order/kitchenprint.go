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
	return SyncPrintKitchenWithDiff(ctx, svcCtx, e, banner, nil)
}

// SyncPrintKitchenWithDiff 同步打印，并附带与旧订单的差异（新增/减量/删除），渲染时给出标记。
func SyncPrintKitchenWithDiff(ctx context.Context, svcCtx *svc.ServiceContext, e *ent.Order, banner string, diff *OrderItemDiff) error {
	if svcCtx == nil || e == nil {
		return errors.New("内部错误")
	}
	if svcCtx.Spyun == nil {
		return errors.New("商鹏云打印未启用或未配置")
	}
	div := kitchenTicketAmountDivisor(svcCtx.Config.Spyun.KitchenTicketAmountScale)
	content := formatKitchenTicket(e, banner, div, diff)
	_, err := svcCtx.Spyun.PrintOrder(ctx, "", content, 1)
	return err
}

// scheduleKitchenPrint 异步提交商鹏打印；nil 客户端或未启用时不做事。失败只记日志，不影响订单接口结果。
func scheduleKitchenPrint(svcCtx *svc.ServiceContext, e *ent.Order, banner string) {
	scheduleKitchenPrintWithDiff(svcCtx, e, banner, nil)
}

// scheduleKitchenPrintWithDiff 异步打印，并附带 diff；diff 为 nil 时与 scheduleKitchenPrint 行为一致。
func scheduleKitchenPrintWithDiff(svcCtx *svc.ServiceContext, e *ent.Order, banner string, diff *OrderItemDiff) {
	if svcCtx == nil || svcCtx.Spyun == nil || e == nil {
		return
	}
	div := kitchenTicketAmountDivisor(svcCtx.Config.Spyun.KitchenTicketAmountScale)
	content := formatKitchenTicket(e, banner, div, diff)
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

// ─────────────────────────── 排版常量 ───────────────────────────

// 58mm 热敏一行约 32 半角列；汉字按 2 列计。
const (
	ticketLineWidth    = 32
	ticketIndexWidth   = 3 // "1. " 序号宽度
	ticketQtyMaxWidth  = 5 // " ×999" 数量列预留宽度
	ticketHeavyRuleCh  = "="
	ticketLightRuleCh  = "-"

	// ticketCurrency 货币符号。
	// 半角 ¥ (U+00A5) 在多数热敏机字体里笔画细、甚至缺字形，打印模糊。
	// 改用全角 ￥ (U+FFE5)：CJK 字体内置、粗体清晰，宽度 2 列与中文对齐自然。
	// 若某型号也不支持，可改为 "元 " 作后缀（如 28.00元）。
	ticketCurrency = "￥"
)

// ─────────────────────────── 主流程 ───────────────────────────

// formatKitchenTicket 渲染厨房单完整内容。diff 为 nil 时按普通新单格式打印（无新增/删除标记）。
func formatKitchenTicket(e *ent.Order, banner string, amountDivisor int64, diff *OrderItemDiff) string {
	var b strings.Builder

	// 顶部 banner（新单 / 加菜 / 补打）
	b.WriteString(renderBannerBlock(banner))

	// 桌台 / 外带（大字醒目）
	b.WriteString(renderHeadlineBlock(e))

	// 订单号与时间（小字）
	b.WriteString(renderMetaBlock(e))

	// 当前菜品列表（带 [新加]/[改量] 标记）
	b.WriteString(heavyRuleLine())
	items, _ := e.Edges.ItemsOrErr()
	totalQty := 0
	for idx, it := range items {
		if it == nil {
			continue
		}
		b.WriteString(renderItemBlock(idx+1, it, amountDivisor, diff))
		totalQty += it.Quantity
	}
	b.WriteString(heavyRuleLine())

	// 汇总（道数/份数 + 合计/实收）
	b.WriteString(renderTotalsBlock(len(items), totalQty, e, amountDivisor))

	// 整单备注（加粗显眼）
	if e.Remark != nil && strings.TrimSpace(*e.Remark) != "" {
		b.WriteString(renderRemarkBlock(strings.TrimSpace(*e.Remark)))
	}

	// 已删除菜（带删除线，单独成区）
	if diff != nil && len(diff.Removed) > 0 {
		b.WriteString(renderRemovedBlock(diff.Removed, amountDivisor))
	}

	b.WriteString("<CUT>")
	return b.String()
}

// ─────────────────────────── 顶部 banner ───────────────────────────

// renderBannerBlock 按 banner 含义渲染醒目横幅；空 banner 时不输出。
// 已识别语义：[新单] / [改单重打] / [手动打印]，其他原样显示。
//
// 排版：加粗居中 + 倍宽（不倍高，不加上下分隔线），避免占用双行高度与额外空行。
// 倍宽 <W> 让字横向放大、纵向仍是 1 行；醒目度足够区分新单/加菜/补打。
func renderBannerBlock(banner string) string {
	banner = strings.TrimSpace(banner)
	if banner == "" {
		return ""
	}
	title, warn := bannerTitleAndWarning(banner)

	var b strings.Builder
	b.WriteString("<C><B><W>")
	b.WriteString(escapeSpyunText(title))
	b.WriteString("</W></B></C><BR>")
	if warn != "" {
		b.WriteString("<C><B>")
		b.WriteString(escapeSpyunText(warn))
		b.WriteString("</B></C><BR>")
	}
	return b.String()
}

func bannerTitleAndWarning(banner string) (title, warn string) {
	// 去掉首尾的方括号，便于匹配中文
	trimmed := strings.Trim(banner, "[]【】 ")
	switch trimmed {
	case "新单":
		return "★ 新  单 ★", ""
	case "改单重打":
		// 现行实现是 ReplaceItems，整单重打。提醒厨房核对已制作。
		return "★ 加  菜 ★", "※ 全单重打 请核对已做 ※"
	case "手动打印":
		return "★ 补  打 ★", ""
	default:
		return banner, ""
	}
}

// ─────────────────────────── 桌台行 ───────────────────────────

func renderHeadlineBlock(e *ent.Order) string {
	var b strings.Builder
	switch string(e.OrderType) {
	case "dine_in":
		tbl, err := e.Edges.TableOrErr()
		if err == nil && tbl != nil {
			line := "桌台: " + strings.TrimSpace(tbl.Code)
			cat, _ := tbl.Edges.CategoryOrErr()
			if cat != nil && strings.TrimSpace(cat.Name) != "" {
				line += " (" + strings.TrimSpace(cat.Name) + ")"
			}
			// 加粗倍高，桌台是厨师除菜名外最关心的字段
			b.WriteString("<B><H>")
			b.WriteString(escapeSpyunText(line))
			b.WriteString("</H></B><BR>")
		} else {
			b.WriteString("<B><H>桌台: 堂食</H></B><BR>")
		}
	case "takeaway":
		b.WriteString("<B><H>※ 外  带 ※</H></B><BR>")
	default:
		b.WriteString("<B><H>类型: ")
		b.WriteString(escapeSpyunText(string(e.OrderType)))
		b.WriteString("</H></B><BR>")
	}
	return b.String()
}

// ─────────────────────────── 元信息（订单号/时间）───────────────────────────

func renderMetaBlock(e *ent.Order) string {
	var b strings.Builder
	// 下单/打印时间：厨师只关心今天的时辰，只显示 HH:MM
	createdAt := e.CreatedAt.Format("15:04")
	printedAt := time.Now().Format("15:04")
	b.WriteString(fmt.Sprintf("下单: %s    打印: %s<BR>", createdAt, printedAt))
	b.WriteString("订单: ")
	b.WriteString(escapeSpyunText(e.OrderNo))
	b.WriteString("<BR>")
	return b.String()
}

// ─────────────────────────── 菜品块 ───────────────────────────

// renderItemBlock 单个菜品的多行块。
// 主行：序号. [新加]菜名              ×数量
// 副行：   ￥单价 × 数量 = ￥金额
// 改量：   ※ 改量: 原×3 → 现×1 ※
// 规格：   规格: xxx
func renderItemBlock(idx int, it *ent.OrderItem, div int64, diff *OrderItemDiff) string {
	var b strings.Builder

	indexStr := fmt.Sprintf("%d.", idx)
	indexPad := padASCIIRight(indexStr, ticketIndexWidth)

	// 判定变更类型（仅 diff 非空时）
	kind, oldQty := ItemDiffNone, 0
	if diff != nil {
		if d, ok := diff.ByKey[itemKey(it)]; ok {
			kind = d.Kind
			oldQty = d.OldQty
		}
	}

	// 菜名前缀（新增标记），最显眼，加粗
	var namePrefix string
	switch kind {
	case ItemDiffAdded:
		namePrefix = "[新加] "
	}

	qtyStr := fmt.Sprintf("×%d", it.Quantity)
	qtyWidth := ticketDisplayWidth(qtyStr)
	if qtyWidth > ticketQtyMaxWidth {
		qtyWidth = ticketQtyMaxWidth
	}

	nameMax := ticketLineWidth - ticketIndexWidth - qtyWidth - 1
	if nameMax < 4 {
		nameMax = 4
	}
	rawName := escapeSpyunText(it.MenuName)
	if namePrefix != "" {
		rawName = namePrefix + rawName
	}
	name := truncateTicketNameDisplay(rawName, nameMax)
	leftPart := indexPad + padDisplayRight(name, ticketLineWidth-ticketIndexWidth-qtyWidth)

	b.WriteString(leftPart)
	b.WriteString("<B>")
	b.WriteString(qtyStr)
	b.WriteString("</B><BR>")

	// 数量变化：单独一行警示，加粗
	if kind == ItemDiffQtyChanged {
		b.WriteString("   <B>※ 改量: 原×")
		b.WriteString(fmt.Sprintf("%d → 现×%d ※</B><BR>", oldQty, it.Quantity))
	}

	// 副行：价格细节
	amtRaw := orderLineStoredAmount(it)
	priceLine := fmt.Sprintf("   %s%s × %d = %s%s",
		ticketCurrency, fmtTicketMoney(it.UnitPrice, div),
		it.Quantity,
		ticketCurrency, fmtTicketMoney(amtRaw, div),
	)
	b.WriteString(escapeSpyunText(priceLine))
	b.WriteString("<BR>")

	// 规格
	if it.SpecInfo != nil {
		spec := strings.TrimSpace(*it.SpecInfo)
		if spec != "" {
			b.WriteString("   规格: ")
			b.WriteString(escapeSpyunText(spec))
			b.WriteString("<BR>")
		}
	}
	return b.String()
}

// ─────────────────────────── 已删除菜块 ───────────────────────────

// strikeStyle 删除线渲染样式：
//
//	"underline" 下方紧贴一行 ━ 模拟（最稳，所有热敏机都支持）
//	"unicode"   字符后插 U+0336 组合长横线（视觉最接近"字上有横线"，但依赖打印机字体支持）
//
// 默认 underline；实测 unicode 在自家打印机上可用时可切换。
const strikeStyle = "underline"

// renderRemovedBlock 渲染已删除菜列表。
// 单道菜本身已带 [删] 前缀 + 下方横线作为视觉标识，因此不再额外加大标题块；
// 仅用一条细分隔线把"已删除"和上方汇总/备注隔开。
func renderRemovedBlock(removed []*ent.OrderItem, div int64) string {
	if len(removed) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString(lightRuleLine())
	for _, it := range removed {
		if it == nil {
			continue
		}
		b.WriteString(renderRemovedItemBlock(it, div))
	}
	return b.String()
}

// renderRemovedItemBlock 单道已删菜品的多行块（带删除线视觉）。
//
//	[删] 宫保鸡丁                ×2
//	     ━━━━━━━━━━━━━━━━━━━━━━━━━━━
//	     原: ￥28.00 × 2 = ￥56.00
func renderRemovedItemBlock(it *ent.OrderItem, div int64) string {
	var b strings.Builder

	prefix := "[删] "
	qtyStr := fmt.Sprintf("×%d", it.Quantity)
	qtyWidth := ticketDisplayWidth(qtyStr)
	if qtyWidth > ticketQtyMaxWidth {
		qtyWidth = ticketQtyMaxWidth
	}
	nameMax := ticketLineWidth - qtyWidth - 1
	rawName := prefix + escapeSpyunText(it.MenuName)
	name := truncateTicketNameDisplay(rawName, nameMax)

	if strikeStyle == "unicode" {
		// 字符后插 U+0336 组合长横线，整体仍按显示宽度对齐
		struck := applyUnicodeStrike(name)
		leftPart := padDisplayRight(struck, ticketLineWidth-qtyWidth)
		b.WriteString(leftPart)
		b.WriteString(qtyStr)
		b.WriteString("<BR>")
	} else {
		// underline：菜名正常一行，下方紧贴一行 ━ 模拟划掉
		leftPart := padDisplayRight(name, ticketLineWidth-qtyWidth)
		b.WriteString(leftPart)
		b.WriteString(qtyStr)
		b.WriteString("<BR>")
		// 下方横线，长度按菜名显示宽度
		underline := strings.Repeat("━", (ticketDisplayWidth(name)+1)/2)
		b.WriteString(underline)
		b.WriteString("<BR>")
	}

	// 副行：原价格
	amtRaw := orderLineStoredAmount(it)
	priceLine := fmt.Sprintf("     原: %s%s × %d = %s%s",
		ticketCurrency, fmtTicketMoney(it.UnitPrice, div),
		it.Quantity,
		ticketCurrency, fmtTicketMoney(amtRaw, div),
	)
	b.WriteString(escapeSpyunText(priceLine))
	b.WriteString("<BR>")

	if it.SpecInfo != nil {
		spec := strings.TrimSpace(*it.SpecInfo)
		if spec != "" {
			b.WriteString("     规格: ")
			b.WriteString(escapeSpyunText(spec))
			b.WriteString("<BR>")
		}
	}
	return b.String()
}

// applyUnicodeStrike 给每个 rune 后追加 U+0336，模拟"字上横线"。
// 注意：能否正确显示取决于热敏机字体是否支持组合标记叠加。
func applyUnicodeStrike(s string) string {
	var b strings.Builder
	for _, r := range s {
		b.WriteRune(r)
		b.WriteRune('\u0336')
	}
	return b.String()
}

// ─────────────────────────── 汇总块 ───────────────────────────

func renderTotalsBlock(courses, totalQty int, e *ent.Order, div int64) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("共 %d 道 / %d 份<BR>", courses, totalQty))
	b.WriteString(fmt.Sprintf("合计: %s%s    实收: %s%s<BR>",
		ticketCurrency, fmtTicketMoney(e.TotalAmount, div),
		ticketCurrency, fmtTicketMoney(e.ActualAmount, div),
	))
	return b.String()
}

// ─────────────────────────── 备注块 ───────────────────────────

func renderRemarkBlock(remark string) string {
	var b strings.Builder
	b.WriteString(lightRuleLine())
	b.WriteString("<B>★ 备注: ")
	b.WriteString(escapeSpyunText(remark))
	b.WriteString(" ★</B><BR>")
	return b.String()
}

// ─────────────────────────── 分隔线 ───────────────────────────

func heavyRuleLine() string {
	return strings.Repeat(ticketHeavyRuleCh, ticketLineWidth) + "<BR>"
}

func lightRuleLine() string {
	return strings.Repeat(ticketLightRuleCh, ticketLineWidth) + "<BR>"
}

// 兼容旧调用（测试中仍有引用）：以前的 "----" 分隔线。
func kitchenTicketRuleLine() string {
	return lightRuleLine()
}

// ─────────────────────────── 文本宽度工具 ───────────────────────────

// ticketDisplayWidth 估算热敏一行占用列数（与常见 58mm 汉字倍宽规则一致）。
// ASCII 占 1 列，其余 Unicode 占 2 列。
func ticketDisplayWidth(s string) int {
	w := 0
	for _, r := range s {
		if r <= 0x007F {
			w++
		} else {
			w += 2
		}
	}
	return w
}

// padDisplayRight 右侧补半角空格，使显示列数达到 target（左对齐）。
func padDisplayRight(s string, target int) string {
	for ticketDisplayWidth(s) < target {
		s += " "
	}
	return s
}

// padASCIIRight 仅按 ASCII 长度右侧补齐。
func padASCIIRight(s string, width int) string {
	for len(s) < width {
		s += " "
	}
	return s
}

// truncateTicketNameDisplay 按显示列数截断菜名，过长加 …。
func truncateTicketNameDisplay(s string, maxDisplay int) string {
	s = strings.TrimSpace(s)
	if maxDisplay <= 0 {
		return ""
	}
	if ticketDisplayWidth(s) <= maxDisplay {
		return s
	}
	const ell = "…"
	budget := maxDisplay - ticketDisplayWidth(ell)
	if budget < 1 {
		return ell
	}
	var b strings.Builder
	for _, ch := range s {
		trial := b.String() + string(ch)
		if ticketDisplayWidth(trial) > budget {
			break
		}
		b.WriteRune(ch)
	}
	out := b.String()
	if out == "" {
		return ell
	}
	return out + ell
}

// ─────────────────────────── 金额工具 ───────────────────────────

// orderLineStoredAmount 行金额：优先库内小计，否则单价×数量。
func orderLineStoredAmount(it *ent.OrderItem) int64 {
	if it == nil {
		return 0
	}
	if it.Amount != 0 {
		return it.Amount
	}
	return it.UnitPrice * int64(it.Quantity)
}

// kitchenTicketAmountDivisor 根据配置：cent/fen/cents → 分转元(/100)；否则按元直出。
func kitchenTicketAmountDivisor(scale string) int64 {
	s := strings.ToLower(strings.TrimSpace(scale))
	switch s {
	case "cent", "fen", "cents":
		return 100
	default:
		return 1
	}
}

// fmtTicketMoney 按 divisor 将库内整型格式化为两位小数的「元」展示（divisor=100 时表示库内为分）。
func fmtTicketMoney(amount int64, divisor int64) string {
	if divisor <= 1 {
		return fmt.Sprintf("%.2f", float64(amount))
	}
	return fmt.Sprintf("%.2f", float64(amount)/float64(divisor))
}

// 商鹏内容为文本指令混排，避免顾客备注等破坏标签：去掉尖括号与换行。
func escapeSpyunText(s string) string {
	s = strings.ReplaceAll(s, "<", "＜")
	s = strings.ReplaceAll(s, ">", "＞")
	s = strings.ReplaceAll(s, "\r", "")
	s = strings.ReplaceAll(s, "\n", " ")
	return strings.TrimSpace(s)
}
