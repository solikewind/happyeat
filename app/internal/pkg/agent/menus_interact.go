package agent

import (
	"context"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/openai"
)

type MenusTechAgent struct {
	Agent *blades.Agent
}

func NewMenusTechAgent(c *Config) (*MenusTechAgent, error) {
	model := openai.NewModel("qwen3.5-flash", openai.Config{
		APIKey:  c.APIKey,
		BaseURL: c.BaseURL,
	})

	agent, err := blades.NewAgent("MenusTechAgent",
		blades.WithModel(model),
		blades.WithInstruction("你是一个菜单技术专家，专门为 HappyEat 餐厅管理系统提供服务。你可以帮助用户处理菜单、订单等相关问题。"),
	)
	if err != nil {
		return nil, err
	}

	return &MenusTechAgent{
		Agent: &agent,
	}, nil
}

// 1 按拼音搜索菜单 2 前端添加菜单到订单列表中
func (a *MenusTechAgent) SearchMenus(ctx context.Context, prompt string) (string, error) {
	input := blades.UserMessage(prompt)

	runner := blades.NewRunner(*a.Agent)
	result, err := runner.Run(ctx, input)
	if err != nil {
		return "", err
	}
	return result.Text, nil
}
