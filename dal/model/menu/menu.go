// Package menu 提供菜单与分类的 data 逻辑。
package menu

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/solikewind/happyeat/dal/model/ent"
	entmenu "github.com/solikewind/happyeat/dal/model/ent/menu"
	"github.com/solikewind/happyeat/dal/model/ent/menucategory"
	"github.com/solikewind/happyeat/dal/model/ent/menuspec"
)

// Menu 菜单数据访问。
type Menu struct {
	c *ent.Client
}

// NewMenu 创建 Menu。
func NewMenu(c *ent.Client) *Menu {
	return &Menu{c: c}
}

// SpecInput 规格项入参。
type SpecInput struct {
	SpecType   string
	SpecValue  string
	PriceDelta float64
}

// CreateMenuInput 创建菜单入参。
type CreateMenuInput struct {
	Name        string
	Description string
	Image       string
	Price       float64
	CategoryID  int
	Specs       []SpecInput
}

// Create 创建菜单及规格（事务内）。
func (m *Menu) Create(ctx context.Context, in CreateMenuInput) (*ent.Menu, error) {
	tx, err := m.c.Tx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	create := tx.Menu.Create().
		SetName(in.Name).
		SetPrice(in.Price).
		SetCategoryID(in.CategoryID)
	if in.Description != "" {
		create = create.SetDescription(in.Description)
	}
	if in.Image != "" {
		create = create.SetImage(in.Image)
	}
	entMenu, err := create.Save(ctx)
	if err != nil {
		return nil, err
	}

	for i, spec := range in.Specs {
		_, err = tx.MenuSpec.Create().
			SetMenuID(entMenu.ID).
			SetSpecType(spec.SpecType).
			SetSpecValue(spec.SpecValue).
			SetPriceDelta(spec.PriceDelta).
			SetSort(i).
			Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return m.c.Menu.Query().Where(entmenu.IDEQ(entMenu.ID)).WithCategory().WithSpecs().Only(ctx)
}

// GetByID 按 ID 获取菜单（含 category、specs）。
func (m *Menu) GetByID(ctx context.Context, id int) (*ent.Menu, error) {
	return m.c.Menu.Query().
		Where(entmenu.IDEQ(id)).
		WithCategory().
		WithSpecs().
		Only(ctx)
}

// ListMenusFilter 列表筛选。
type ListMenusFilter struct {
	Name         string
	CategoryName string
	Offset       int
	Limit        int
}

// List 分页列出菜单（含 category、specs），返回列表与总数。
// 统一支持拼音搜索：无论输入中文还是拼音，都使用统一的搜索逻辑。
func (m *Menu) List(ctx context.Context, f ListMenusFilter) ([]*ent.Menu, int64, error) {
	q := m.c.Menu.Query().WithCategory().WithSpecs()

	// 分类筛选在数据库层面处理
	if f.CategoryName != "" {
		q = q.Where(entmenu.HasCategoryWith(menucategory.NameEQ(f.CategoryName)))
	}

	var allList []*ent.Menu
	var err error

	// 如果有关键词，统一使用拼音匹配（支持中文和拼音）
	if f.Name != "" {
		// 先查询所有符合分类条件的菜单（不分页）
		allList, err = q.Order(entmenu.ByID(sql.OrderDesc())).All(ctx)
		if err != nil {
			return nil, 0, err
		}

		// 在内存中做拼音匹配（统一处理中文和拼音）
		matchedList := make([]*ent.Menu, 0)
		for _, menu := range allList {
			if MatchPinyin(menu.Name, f.Name) {
				matchedList = append(matchedList, menu)
			}
		}
		allList = matchedList
	} else {
		// 没有关键词，直接查询所有
		allList, err = q.Order(entmenu.ByID(sql.OrderDesc())).All(ctx)
		if err != nil {
			return nil, 0, err
		}
	}

	// 计算总数
	total := int64(len(allList))

	// 手动分页
	if f.Limit <= 0 {
		f.Limit = 10
	}

	start := f.Offset
	end := start + f.Limit
	if start > len(allList) {
		start = len(allList)
	}
	if end > len(allList) {
		end = len(allList)
	}

	if start >= end {
		return []*ent.Menu{}, total, nil
	}

	list := allList[start:end]
	return list, total, nil
}

// UpdateMenuInput 更新菜单入参。
type UpdateMenuInput struct {
	Name        string
	Description string
	Image       string
	Price       float64
	CategoryID  int
	Specs       []SpecInput
}

// Update 更新菜单及规格（先删后建 specs，事务内）。
func (m *Menu) Update(ctx context.Context, id int, in UpdateMenuInput) error {
	tx, err := m.c.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	upd := tx.Menu.UpdateOneID(id).
		SetName(in.Name).
		SetPrice(in.Price).
		SetCategoryID(in.CategoryID)
	if in.Description != "" {
		upd = upd.SetDescription(in.Description)
	} else {
		upd = upd.ClearDescription()
	}
	if in.Image != "" {
		upd = upd.SetImage(in.Image)
	} else {
		upd = upd.ClearImage()
	}
	if _, err = upd.Save(ctx); err != nil {
		return err
	}

	_, err = tx.MenuSpec.Delete().Where(menuspec.HasMenuWith(entmenu.IDEQ(id))).Exec(ctx)
	if err != nil {
		return err
	}

	for i, spec := range in.Specs {
		_, err = tx.MenuSpec.Create().
			SetMenuID(id).
			SetSpecType(spec.SpecType).
			SetSpecValue(spec.SpecValue).
			SetPriceDelta(spec.PriceDelta).
			SetSort(i).
			Save(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

// Delete 删除菜单及其规格（事务内）。
func (m *Menu) Delete(ctx context.Context, id int) error {
	tx, err := m.c.Tx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.MenuSpec.Delete().Where(menuspec.HasMenuWith(entmenu.IDEQ(id))).Exec(ctx)
	if err != nil {
		return err
	}
	if err = tx.Menu.DeleteOneID(id).Exec(ctx); err != nil {
		return err
	}

	return tx.Commit()
}
