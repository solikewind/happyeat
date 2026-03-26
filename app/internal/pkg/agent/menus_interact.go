package agent

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-kratos/blades"
	"github.com/go-kratos/blades/contrib/openai"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/menu"
)

const (
	menuLookupLimit = 10
	menusModelName  = "qwen3.5-flash"
)

var (
	fallbackQuantityMap = map[string]int{
		"一": 1, "二": 2, "两": 2, "三": 3, "四": 4, "五": 5,
		"六": 6, "七": 7, "八": 8, "九": 9, "十": 10,
		"一份": 1, "两份": 2, "两分": 2, "三分": 3, "四分": 4, "五分": 5,
		"六分": 6, "七分": 7, "八分": 8, "九分": 9, "十分": 10,
		"一个": 1, "两个": 2, "三个": 3, "四个": 4, "五个": 5,
		"六个": 6, "七个": 7, "八个": 8, "九个": 9, "十个": 10,
	}
	fallbackPromptPatterns = []*regexp.Regexp{
		regexp.MustCompile(`([一二两三四五六七八九十\d]+)[份个]?([^，。,.\n]+)`),
		regexp.MustCompile(`(\d+)[份个]?([^，。,.\n]+)`),
	}
	menuNameLeadingPattern  = regexp.MustCompile(`^(来|要|点|给|来一份|来一个|要一份|要一个|点一份|点一个)`)
	menuNameTrailingPattern = regexp.MustCompile(`(一份|一个|份|个|菜)?$`)
)

type MenuItem struct {
	Menu     *ent.Menu
	Quantity int
	MenuName string
}

type MenusTechAgent struct {
	Agent *blades.Agent
	Menu  *menu.Menu
}

type extractedMenuItem struct {
	MenuName string `json:"menu_name"`
	Quantity int    `json:"quantity"`
}

func NewMenusTechAgent(c *Config, menuData *menu.Menu) (*MenusTechAgent, error) {
	modelName := menusModelName
	if strings.TrimSpace(c.Model) != "" {
		modelName = strings.TrimSpace(c.Model)
	}

	model := openai.NewModel(modelName, openai.Config{
		APIKey:  c.APIKey,
		BaseURL: c.BaseURL,
	})

	agent, err := blades.NewAgent("MenusTechAgent",
		blades.WithModel(model),
		blades.WithInstruction("你是一个菜单助手，负责从用户点餐描述中提取菜名和数量。"),
	)
	if err != nil {
		return nil, err
	}

	return &MenusTechAgent{
		Agent: &agent,
		Menu:  menuData,
	}, nil
}

func (a *MenusTechAgent) ParseOrderPrompt(ctx context.Context, prompt string) ([]MenuItem, error) {
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return nil, nil
	}
	if a.Menu == nil {
		return nil, errors.New("menu data source is not initialized")
	}
	if a.Agent == nil {
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	extractPrompt := fmt.Sprintf(`请从以下文本中提取所有菜单名称和对应数量。
文本：%s

请只返回 JSON 数组，格式如下：
[
  {"menu_name": "菜单名称", "quantity": 数量}
]

如果数量未明确说明，默认为 1。`, prompt)

	input := blades.UserMessage(extractPrompt)
	runner := blades.NewRunner(*a.Agent)
	result, err := runner.Run(ctx, input)
	if err != nil {
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	jsonStr := extractJSONFromText(result.Text())
	if jsonStr == "" {
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	var parsed []extractedMenuItem
	if err := json.Unmarshal([]byte(jsonStr), &parsed); err != nil {
		return a.parseOrderPromptFallback(ctx, prompt)
	}

	return a.buildMenuItems(ctx, parsed), nil
}

func (a *MenusTechAgent) parseOrderPromptFallback(ctx context.Context, prompt string) ([]MenuItem, error) {
	if a.Menu == nil {
		return nil, errors.New("menu data source is not initialized")
	}

	resultItems := make([]MenuItem, 0)
	processed := make(map[string]bool)

	for _, pattern := range fallbackPromptPatterns {
		matches := pattern.FindAllStringSubmatch(prompt, -1)
		for _, match := range matches {
			if len(match) < 3 {
				continue
			}

			quantityStr := strings.TrimSpace(match[1])
			menuName := strings.TrimSpace(match[2])
			menuName = menuNameLeadingPattern.ReplaceAllString(menuName, "")
			menuName = menuNameTrailingPattern.ReplaceAllString(menuName, "")
			menuName = strings.TrimSpace(menuName)

			if menuName == "" || processed[menuName] {
				continue
			}
			processed[menuName] = true

			quantity := 1
			if q, ok := fallbackQuantityMap[quantityStr]; ok {
				quantity = q
			} else if q, err := strconv.Atoi(quantityStr); err == nil {
				quantity = q
			}

			matchedMenu, err := a.findMenuByName(ctx, menuName)
			if err != nil || matchedMenu == nil {
				continue
			}

			resultItems = append(resultItems, MenuItem{
				Menu:     matchedMenu,
				Quantity: quantity,
				MenuName: menuName,
			})
		}
	}

	return resultItems, nil
}

func extractJSONFromText(text string) string {
	start := strings.Index(text, "[")
	if start == -1 {
		return ""
	}

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

func (a *MenusTechAgent) SearchMenusWithPrompt(ctx context.Context, prompt string) (string, error) {
	items, err := a.ParseOrderPrompt(ctx, prompt)
	if err != nil {
		return "", fmt.Errorf("解析提示词失败: %w", err)
	}

	if len(items) == 0 {
		return "未找到匹配的菜单", nil
	}

	resultText := fmt.Sprintf("找到 %d 个菜单：\n\n", len(items))
	totalPrice := int64(0)
	for i, item := range items {
		price := item.Menu.Price * int64(item.Quantity)
		totalPrice += price
		resultText += fmt.Sprintf("%d. %s × %d = ¥%.2f\n", i+1, item.Menu.Name, item.Quantity, float64(price)/100.0)
		if item.Menu.Description != nil && *item.Menu.Description != "" {
			resultText += fmt.Sprintf("   描述: %s\n", *item.Menu.Description)
		}
		if item.Menu.Edges.Category != nil {
			resultText += fmt.Sprintf("   分类: %s\n", item.Menu.Edges.Category.Name)
		}
		resultText += "\n"
	}
	resultText += fmt.Sprintf("总计: ¥%.2f\n", float64(totalPrice)/100.0)

	return resultText, nil
}

func (a *MenusTechAgent) buildMenuItems(ctx context.Context, extracted []extractedMenuItem) []MenuItem {
	result := make([]MenuItem, 0, len(extracted))
	processed := make(map[string]bool)

	for _, item := range extracted {
		menuName := strings.TrimSpace(item.MenuName)
		if menuName == "" || processed[menuName] {
			continue
		}
		processed[menuName] = true

		quantity := item.Quantity
		if quantity <= 0 {
			quantity = 1
		}

		matchedMenu, err := a.findMenuByName(ctx, menuName)
		if err != nil || matchedMenu == nil {
			continue
		}

		result = append(result, MenuItem{
			Menu:     matchedMenu,
			Quantity: quantity,
			MenuName: menuName,
		})
	}

	return result
}

func (a *MenusTechAgent) findMenuByName(ctx context.Context, menuName string) (*ent.Menu, error) {
	list, _, err := a.Menu.List(ctx, menu.ListMenusFilter{
		Name:         menuName,
		CategoryName: "",
		Offset:       0,
		Limit:        menuLookupLimit,
	})
	if err != nil {
		return nil, err
	}
	if len(list) == 0 {
		return nil, nil
	}
	return list[0], nil
}
