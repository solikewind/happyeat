package agent

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/openai"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/menu"
)

// MenuItem 菜单项（包含菜单信息和数量）
type MenuItem struct {
	Menu     *ent.Menu // 匹配到的菜单
	Quantity int       // 数量
	MenuName string    // 原始菜单名称（用于显示）
}

type MenusTechAgent struct {
	Agent *blades.Agent
	Menu  *menu.Menu // 菜单数据访问层
}

func NewMenusTechAgent(c *Config, menuData *menu.Menu) (*MenusTechAgent, error) {
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
		Menu:  menuData,
	}, nil
}

// ParseOrderPrompt 从提示词中提取所有菜单及其份数
// 例如："我要两份宫保鸡丁，一个红烧肉，三份麻婆豆腐"
// 返回：[]MenuItem，包含每个菜单的匹配结果和数量
func (a *MenusTechAgent) ParseOrderPrompt(ctx context.Context, prompt string) ([]MenuItem, error) {
	// 如果 Agent 为 nil，直接使用备用方法
	if a.Agent == nil {
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	// 使用 LLM 解析提示词，提取菜单和数量
	extractPrompt := fmt.Sprintf(`请从以下文本中提取所有菜单名称和对应的数量。
文本：%s

请以 JSON 格式返回结果，格式如下：
[
  {"menu_name": "菜单名称", "quantity": 数量},
  ...
]

只返回 JSON 数组，不要其他文字说明。如果数量未明确说明，默认为1。`, prompt)

	// 调用 LLM 提取菜单信息
	input := blades.UserMessage(extractPrompt)
	runner := blades.NewRunner(*a.Agent)
	result, err := runner.Run(ctx, input)
	if err != nil {
		// 如果 LLM 失败，使用简单的正则表达式解析
		fmt.Println("ParseOrderPrompt error:", err)
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	// 解析 LLM 返回的 JSON
	var menuItems []struct {
		MenuName string `json:"menu_name"`
		Quantity int    `json:"quantity"`
	}

	// 尝试从结果中提取 JSON（可能包含其他文本）
	// result 是 *blades.Message 类型，需要获取其文本内容
	resultText := result.Text()
	jsonStr := extractJSONFromText(resultText)
	if jsonStr == "" {
		// 如果提取不到 JSON，使用备用方法
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	if err := json.Unmarshal([]byte(jsonStr), &menuItems); err != nil {
		// JSON 解析失败，使用备用方法
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	// 对每个菜单名称进行搜索匹配
	resultItems := make([]MenuItem, 0)
	for _, item := range menuItems {
		if item.MenuName == "" {
			continue
		}
		if item.Quantity <= 0 {
			item.Quantity = 1
		}

		// 搜索匹配的菜单
		list, _, err := a.Menu.List(ctx, menu.ListMenusFilter{
			Name:         item.MenuName,
			CategoryName: "",
			Offset:       0,
			Limit:        10, // 每个菜单最多返回10个匹配结果
		})
		if err != nil {
			continue
		}

		// 取第一个匹配结果
		if len(list) > 0 {
			resultItems = append(resultItems, MenuItem{
				Menu:     list[0],
				Quantity: item.Quantity,
				MenuName: item.MenuName,
			})
		}
	}

	return resultItems, nil
}

// parseOrderPromptFallback 备用解析方法（使用正则表达式）
func (a *MenusTechAgent) parseOrderPromptFallback(ctx context.Context, prompt string) ([]MenuItem, error) {
	// 数量映射表
	quantityMap := map[string]int{
		"一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5,
		"六": 6, "七": 7, "八": 8, "九": 9, "十": 10,
		"一份": 1, "两份": 2, "两分": 2, "三分": 3, "四分": 4, "五分": 5,
		"六分": 6, "七分": 7, "八分": 8, "九分": 9, "十分": 10,
		"一个": 1, "两个": 2, "三个": 3, "四个": 4, "五个": 5,
		"六个": 6, "七个": 7, "八个": 8, "九个": 9, "十个": 10,
	}

	// 匹配模式：数量 + 菜单名称
	// 例如：两份宫保鸡丁、一个红烧肉、3份麻婆豆腐
	patterns := []*regexp.Regexp{
		regexp.MustCompile(`([一二两三五六七八九十\d]+)[份个]?([^，,。.\n]+)`),
		regexp.MustCompile(`(\d+)[份个]?([^，,。.\n]+)`),
	}

	resultItems := make([]MenuItem, 0)
	processed := make(map[string]bool)

	for _, pattern := range patterns {
		matches := pattern.FindAllStringSubmatch(prompt, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}

			quantityStr := strings.TrimSpace(match[1])
			menuName := strings.TrimSpace(match[2])

			// 移除常见的前缀和后缀
			menuName = regexp.MustCompile(`^(来|要|点|加|上|来一份|来一个|要一份|要一个|点一份|点一个)`).ReplaceAllString(menuName, "")
			menuName = regexp.MustCompile(`(一份|一个|份|个|菜|道)$`).ReplaceAllString(menuName, "")
			menuName = strings.TrimSpace(menuName)

			if menuName == "" || processed[menuName] {
				continue
			}
			processed[menuName] = true

			// 解析数量
			quantity := 1
			if q, ok := quantityMap[quantityStr]; ok {
				quantity = q
			} else if q, err := strconv.Atoi(quantityStr); err == nil {
				quantity = q
			}

			// 搜索匹配的菜单
			list, _, err := a.Menu.List(ctx, menu.ListMenusFilter{
				Name:         menuName,
				CategoryName: "",
				Offset:       0,
				Limit:        10,
			})
			if err != nil || len(list) == 0 {
				continue
			}

			resultItems = append(resultItems, MenuItem{
				Menu:     list[0],
				Quantity: quantity,
				MenuName: menuName,
			})
		}
	}

	return resultItems, nil
}

// extractJSONFromText 从文本中提取 JSON 数组
func extractJSONFromText(text string) string {
	// 查找 JSON 数组的开始和结束
	start := strings.Index(text, "[")
	if start == -1 {
		return ""
	}

	// 从 [ 开始，找到匹配的 ]
	depth := 0
	for i := start; i < len(text); i++ {
		if text[i] == '[' {
			depth++
		} else if text[i] == ']' {
			depth--
			if depth == 0 {
				return text[start : i+1]
			}
		}
	}

	return ""
}

// SearchMenusWithPrompt 使用 LLM 处理用户提示并搜索菜单
// 现在支持从提示词中提取多个菜单及其份数
func (a *MenusTechAgent) SearchMenusWithPrompt(ctx context.Context, prompt string) (string, error) {
	// 解析提示词，提取所有菜单和数量
	items, err := a.ParseOrderPrompt(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("解析提示词失败: %w", err)
	}

	if len(items) == 0 {
		return "未找到匹配的菜单", nil
	}

	// 格式化结果
	resultText := fmt.Sprintf("找到 %d 个菜单：\n\n", len(items))
	totalPrice := 0.0
	for i, item := range items {
		price := item.Menu.Price * float64(item.Quantity)
		totalPrice += price
		resultText += fmt.Sprintf("%d. %s × %d = ¥%.2f\n", i+1, item.Menu.Name, item.Quantity, price)
		if item.Menu.Description != nil && *item.Menu.Description != "" {
			resultText += fmt.Sprintf("   描述: %s\n", *item.Menu.Description)
		}
		if item.Menu.Edges.Category != nil {
			resultText += fmt.Sprintf("   分类: %s\n", item.Menu.Edges.Category.Name)
		}
		resultText += "\n"
	}
	resultText += fmt.Sprintf("总计: ¥%.2f\n", totalPrice)

	return resultText, nil
}
