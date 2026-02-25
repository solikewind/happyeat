package agent

import (
	"os"

	"github.com/solikewind/happyeat/app/internal/config"
)

type Config struct {
	APIKey  string
	BaseURL string
}

// NewAgent 创建并配置 Blades Agent
func NewConfig(c config.LLM) (*Config, error) {
	// 从环境变量或配置中获取 API Key
	apiKey := c.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("OPENAI_API_KEY")
	}

	// 从环境变量或配置中获取 Base URL（可选，用于兼容其他 OpenAI 兼容的 API）
	baseURL := c.BaseURL
	if baseURL == "" {
		baseURL = os.Getenv("OPENAI_BASE_URL")
	}

	// // 创建 Agent 选项
	// opts := []blades.AgentOption{
	// 	blades.WithModel(model),
	// 	blades.WithInstruction(),
	// }
	// // 你是一个智能助手，专门为 HappyEat 餐厅管理系统提供服务。你可以帮助用户处理订单、菜单、餐桌等相关问题。
	// // 如果有工具配置，可以在这里添加
	// // opts = append(opts, blades.WithTools(...))

	// // 创建 Agent
	// bladesAgent, err := blades.NewAgent("HappyEat Assistant", opts...)
	// if err != nil {
	// 	return nil, err
	// }

	return &Config{
		APIKey:  apiKey,
		BaseURL: baseURL,
	}, nil
}
