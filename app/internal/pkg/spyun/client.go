package spyun

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/solikewind/happyeat/app/internal/config"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	defaultBaseURL   = "https://open.spyun.net"
	pathPrinterPrint = "/v1/printer/print"
)

// Client 商鹏云打印 HTTP 客户端（与 JDL-sdk 类似：配置驱动 + 独立签名 + POST）。
type Client struct {
	appID     string
	appSecret string
	baseURL   string
	sn        string
	http      *http.Client
}

// NewClient 从服务配置构造客户端；未启用或缺少 appid/secret 时返回 nil（与 Cos 一致）。
func NewClient(c config.Spyun) *Client {
	if !c.Enabled {
		return nil
	}
	appID := strings.TrimSpace(c.AppID)
	secret := strings.TrimSpace(c.AppSecret)
	if appID == "" || secret == "" {
		return nil
	}
	base := strings.TrimSpace(c.BaseURL)
	if base == "" {
		base = defaultBaseURL
	}
	base = strings.TrimRight(base, "/")
	timeout := time.Duration(c.TimeoutSec) * time.Second
	if c.TimeoutSec <= 0 {
		timeout = 15 * time.Second
	}
	return &Client{
		appID:     appID,
		appSecret: secret,
		baseURL:   base,
		sn:        strings.TrimSpace(c.SN),
		http: &http.Client{
			Timeout: timeout,
		},
	}
}

// PrintOrder 小票打印。sn 为空时使用配置中的默认 SN；times 为打印份数（1–5，超出则钳制）。
func (c *Client) PrintOrder(ctx context.Context, sn, content string, times int) (*PrintReply, error) {
	if c == nil {
		return nil, nil
	}
	deviceSN := strings.TrimSpace(sn)
	if deviceSN == "" {
		deviceSN = c.sn
	}
	if deviceSN == "" {
		return nil, fmt.Errorf("spyun: 打印机 sn 为空")
	}
	if strings.TrimSpace(content) == "" {
		return nil, fmt.Errorf("spyun: 打印内容为空")
	}
	if times < 1 {
		times = 1
	}
	if times > 5 {
		times = 5
	}

	ts := strconv.FormatInt(time.Now().Unix(), 10)
	params := map[string]string{
		"appid":     c.appID,
		"timestamp": ts,
		"sn":        deviceSN,
		"content":   content,
		"times":     strconv.Itoa(times),
	}
	params["sign"] = BuildSign(params, c.appSecret)

	form := url.Values{}
	for k, v := range params {
		form.Set(k, v)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+pathPrinterPrint, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var out PrintReply
	if err := json.Unmarshal(body, &out); err != nil {
		logx.WithContext(ctx).Errorf("spyun print: invalid json: %s", string(body))
		return nil, fmt.Errorf("spyun: 解析响应失败: %w", err)
	}
	if out.ErrorCode != 0 {
		msg := strings.TrimSpace(out.ErrorMsg)
		if msg == "" {
			msg = fmt.Sprintf("errorcode=%d", out.ErrorCode)
		}
		return &out, fmt.Errorf("spyun: %s", msg)
	}
	return &out, nil
}
