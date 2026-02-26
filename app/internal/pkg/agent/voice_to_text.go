package agent

import (
	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/openai"
)

type VoiceToTextAgent struct {
	Agent *blades.Agent
}

func NewVoiceToTextAgent(c *Config) (*VoiceToTextAgent, error) {
	model := openai.NewModel("fun-asr-realtime", openai.Config{
		APIKey:  c.APIKey,
		BaseURL: c.BaseURL,
	})

	agent, err := blades.NewAgent("VoiceToTextAgent",
		blades.WithModel(model),
		blades.WithInstruction("你是一个语音转文本的智能助手，专门为 HappyEat 餐厅管理系统提供服务。你可以帮助用户处理订单、菜单、餐桌等相关问题。"),
	)
	if err != nil {
		return nil, err
	}

	return &VoiceToTextAgent{
		Agent: &agent,
	}, nil
}

// // Process 处理语音转文本
// func (a *VoiceToTextAgent) Process(ctx context.Context, input string) (string, error) {
// 	response, err := a.Agent.Chat(ctx, input)
// 	if err != nil {
// 		return "", err
// 	}
// 	return response.Text, nil
// }
