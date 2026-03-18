package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/solikewind/happyeat/app/internal/config"
	"github.com/tmc/langchaingo/llms"
	lcopenai "github.com/tmc/langchaingo/llms/openai"
)

const defaultLangChainModel = "qwen-plus"

type OrderIntentItem struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Notes    string `json:"notes,omitempty"`
}

type OrderIntent struct {
	Items []OrderIntentItem `json:"items"`
}

type LangChainService struct {
	model llms.Model
	name  string
}

func NewLangChainService(c config.LLM) (*LangChainService, error) {
	apiKey := strings.TrimSpace(c.APIKey)
	if apiKey == "" {
		return nil, errors.New("llm api key is empty")
	}

	modelName := strings.TrimSpace(c.Model)
	if modelName == "" {
		modelName = defaultLangChainModel
	}

	opts := []lcopenai.Option{
		lcopenai.WithToken(apiKey),
		lcopenai.WithModel(modelName),
	}
	if strings.TrimSpace(c.BaseURL) != "" {
		opts = append(opts, lcopenai.WithBaseURL(strings.TrimSpace(c.BaseURL)))
	}

	model, err := lcopenai.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("init langchain model: %w", err)
	}

	return &LangChainService{
		model: model,
		name:  modelName,
	}, nil
}

func (s *LangChainService) ModelName() string {
	return s.name
}

func (s *LangChainService) ExtractOrderIntent(ctx context.Context, userText string) (*OrderIntent, error) {
	text := strings.TrimSpace(userText)
	if text == "" {
		return &OrderIntent{Items: nil}, nil
	}

	prompt := fmt.Sprintf(
		"你是点餐语义解析器。请从用户输入中提取菜品、数量和备注，严格返回 JSON 对象，不要输出任何额外文字。"+
			`格式：{"items":[{"name":"菜名","quantity":1,"notes":"可选备注"}]}。`+
			"数量默认 1。用户输入：%s", text,
	)

	out, err := llms.GenerateFromSinglePrompt(
		ctx,
		s.model,
		prompt,
		llms.WithModel(s.name),
		llms.WithTemperature(0),
		llms.WithJSONMode(),
	)
	if err != nil {
		return nil, fmt.Errorf("extract order intent: %w", err)
	}

	payload := extractJSONObject(out)
	if payload == "" {
		payload = strings.TrimSpace(out)
	}

	var intent OrderIntent
	if err := json.Unmarshal([]byte(payload), &intent); err != nil {
		return nil, fmt.Errorf("decode order intent json: %w", err)
	}

	for i := range intent.Items {
		if intent.Items[i].Quantity <= 0 {
			intent.Items[i].Quantity = 1
		}
		intent.Items[i].Name = strings.TrimSpace(intent.Items[i].Name)
		intent.Items[i].Notes = strings.TrimSpace(intent.Items[i].Notes)
	}

	return &intent, nil
}

func extractJSONObject(text string) string {
	start := strings.Index(text, "{")
	if start == -1 {
		return ""
	}

	depth := 0
	for i := start; i < len(text); i++ {
		switch text[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
	}
	return ""
}
