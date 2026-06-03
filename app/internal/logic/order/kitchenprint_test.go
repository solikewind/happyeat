package order

import (
	"strings"
	"testing"
)

func TestKitchenTicketAmountDivisor(t *testing.T) {
	if kitchenTicketAmountDivisor("") != 1 {
		t.Fatalf("empty -> 1")
	}
	if kitchenTicketAmountDivisor("yuan") != 1 {
		t.Fatalf("yuan -> 1")
	}
	if kitchenTicketAmountDivisor("cent") != 100 {
		t.Fatalf("cent -> 100")
	}
	if kitchenTicketAmountDivisor("  FEN ") != 100 {
		t.Fatalf("fen -> 100")
	}
}

func TestFmtTicketMoney(t *testing.T) {
	if got := fmtTicketMoney(71, 1); got != "71.00" {
		t.Fatalf("71 yuan: got %q", got)
	}
	if got := fmtTicketMoney(7100, 100); got != "71.00" {
		t.Fatalf("7100 cent: got %q", got)
	}
}

func TestTicketDisplayWidth(t *testing.T) {
	if w := ticketDisplayWidth("AB"); w != 2 {
		t.Fatalf("ascii: %d", w)
	}
	if w := ticketDisplayWidth("炒鸡"); w != 4 {
		t.Fatalf("cjk: %d", w)
	}
	if w := ticketDisplayWidth("炒A鸡"); w != 5 {
		t.Fatalf("mixed: %d", w)
	}
}

func TestTruncateTicketNameDisplay(t *testing.T) {
	if got := truncateTicketNameDisplay("炒鸡", 10); got != "炒鸡" {
		t.Fatalf("no truncate: %q", got)
	}
	got := truncateTicketNameDisplay("麻婆豆腐套餐配米饭", 8)
	if ticketDisplayWidth(got) > 8 {
		t.Fatalf("over width: %q (%d)", got, ticketDisplayWidth(got))
	}
	if !strings.HasSuffix(got, "…") {
		t.Fatalf("expect ellipsis: %q", got)
	}
}

func TestBannerTitleAndWarning(t *testing.T) {
	cases := []struct {
		in        string
		wantTitle string
		wantWarn  bool
	}{
		{"[新单]", "★ 新  单 ★", false},
		{"【新单】", "★ 新  单 ★", false},
		{"[改单重打]", "★ 加  菜 ★", true},
		{"[手动打印]", "★ 补  打 ★", false},
		{"[未识别]", "[未识别]", false},
	}
	for _, c := range cases {
		gotTitle, gotWarn := bannerTitleAndWarning(c.in)
		if gotTitle != c.wantTitle {
			t.Errorf("title for %q: want %q got %q", c.in, c.wantTitle, gotTitle)
		}
		if (gotWarn != "") != c.wantWarn {
			t.Errorf("warn presence for %q: want %v got %q", c.in, c.wantWarn, gotWarn)
		}
	}
}

func TestRenderBannerBlock(t *testing.T) {
	if got := renderBannerBlock(""); got != "" {
		t.Fatalf("empty banner should be empty: %q", got)
	}
	got := renderBannerBlock("[改单重打]")
	for _, want := range []string{
		"<C><B><W>",
		"★ 加  菜 ★",
		"全单重打",
		"</W></B></C>",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("missing %q in %q", want, got)
		}
	}
	if strings.Contains(got, "<H>") || strings.Contains(got, "</H>") {
		t.Errorf("banner should not double height anymore: %q", got)
	}
	if strings.Contains(got, strings.Repeat("=", ticketLineWidth)) {
		t.Errorf("banner should no longer carry heavy rule: %q", got)
	}
}

func TestHeavyAndLightRule(t *testing.T) {
	if got := heavyRuleLine(); !strings.HasPrefix(got, strings.Repeat("=", ticketLineWidth)) {
		t.Fatalf("heavy rule: %q", got)
	}
	if got := lightRuleLine(); !strings.HasPrefix(got, strings.Repeat("-", ticketLineWidth)) {
		t.Fatalf("light rule: %q", got)
	}
}

func TestEscapeSpyunText(t *testing.T) {
	got := escapeSpyunText("a<b>c\nd")
	if strings.ContainsAny(got, "<>\n") {
		t.Fatalf("not escaped: %q", got)
	}
}

func TestApplyUnicodeStrike(t *testing.T) {
	got := applyUnicodeStrike("AB")
	if got != "A\u0336B\u0336" {
		t.Fatalf("strike: %q", got)
	}
}

func TestTicketCurrencyIsFullwidth(t *testing.T) {
	// 必须是全角 ￥ (U+FFE5)，不是半角 ¥ (U+00A5)。
	// 半角在多数热敏机字体上打印不清楚或缺字形。
	if ticketCurrency != "\uFFE5" {
		t.Fatalf("currency should be U+FFE5 ￥, got %q (codepoint %U)", ticketCurrency, []rune(ticketCurrency))
	}
	if ticketCurrency == "\u00A5" {
		t.Fatalf("currency must not be halfwidth ¥ (U+00A5)")
	}
}

