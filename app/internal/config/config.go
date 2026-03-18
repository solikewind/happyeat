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
