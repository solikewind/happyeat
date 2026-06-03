// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	SqlConfig SqlConfig
	Auth      Auth
	Casbin    Casbin
	LLM       LLM
	ASR       ASR
	Cos       Cos
	Spyun     Spyun `json:"spyun,optional"`
}
type SqlConfig struct {
	DataSource string
}

type Auth struct {
	AccessSecret string
	AccessExpire int64
}
type Casbin struct {
	Model string // 模型内联字符串（与 yaml 中 casbin.model 一致）；策略从数据库 casbin_rule 表加载
}

type LLM struct {
	APIKey  string // OpenAI API Key（可选，优先使用环境变量 OPENAI_API_KEY）
	BaseURL string // API Base URL（可选，用于兼容其他 OpenAI 兼容的 API）
	Model   string // 模型名称（如 qwen-plus、gpt-4o-mini）
}

type ASR struct {
	APIKey            string // 阿里云百炼 API Key
	BaseURL           string // 百炼服务基地址
	Model             string // ASR 模型名
	TranscribePath    string // 语音识别接口路径
	HotwordCreatePath string // 热词创建接口路径
}

// Cos 对象存储（腾讯云 COS）。YAML 使用 snake_case（bucket_url 等），需 json 标签才能映射到字段。
// optional：未配置或仅 migrate 时不强制填写。
type Cos struct {
	BucketUrl *string `json:"bucket_url,optional"`
	SecretId  *string `json:"secret_id,optional"`
	SecretKey *string `json:"secret_key,optional"`
}

// Spyun 商鹏云打印（https://www.spyun.net.cn/open/index.html）。未启用或未配密钥时客户端为 nil。
type Spyun struct {
	Enabled     bool   `json:"enabled,optional"`
	AppID       string `json:"app_id,optional"`
	AppSecret   string `json:"app_secret,optional"`
	BaseURL     string `json:"base_url,optional"`     // 默认 https://open.spyun.net
	SN          string `json:"sn,optional"`             // 默认打印机编号
	TimeoutSec  int    `json:"timeout_sec,optional"` // HTTP 超时秒数，默认 15
	// KitchenTicketAmountScale 厨房小票：空或 yuan=库内即元；cent/fen/cents=库内为分（打印时 /100）。
	KitchenTicketAmountScale string `json:"kitchen_ticket_amount_scale,optional"`
}
