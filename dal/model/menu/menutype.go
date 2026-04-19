// Package menu。menutype.go 封装菜单分类相关 data 逻辑。
package menu

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"github.com/solikewind/happyeat/dal/model/ent"
	entmenu "github.com/solikewind/happyeat/dal/model/ent/menu"
	"github.com/solikewind/happyeat/dal/model/ent/menucategory"
)

// MenuType 菜单分类数据访问。
type MenuType struct {
	c *ent.Client
}

// NewMenuType 创建 MenuType。
func NewMenuType(c *ent.Client) *MenuType {
	return &MenuType{c: c}
}

// CreateMenuCategoryInput 创建分类入参。
type CreateMenuCategoryInput struct {
	Name        string
	Description string
	Sort        uint32
}

// Create 创建菜单分类。
func (mt *MenuType) Create(ctx context.Context, in CreateMenuCategoryInput) (*ent.MenuCategory, error) {
	create := mt.c.MenuCategory.Create().SetName(in.Name).SetSort(in.Sort)
	if in.Description != "" {
		create = create.SetDescription(in.Description)
	}

	return create.Save(ctx)
}

// GetByID 按 ID 获取分类。
func (mt *MenuType) GetByID(ctx context.Context, id uint64) (*ent.MenuCategory, error) {
	return mt.c.MenuCategory.Get(ctx, id)
}

// Exist 判断分类是否存在。
func (mt *MenuType) Exist(ctx context.Context, id uint64) (bool, error) {
	return mt.c.MenuCategory.Query().Where(menucategory.IDEQ(id)).Exist(ctx)
}

// ListMenuCategoriesFilter 分类列表筛选。
type ListMenuCategoriesFilter struct {
	Name   string
	Offset int
	Limit  int
}

// List 分页列出分类，返回列表与总数。
func (mt *MenuType) List(ctx context.Context, f ListMenuCategoriesFilter) ([]*ent.MenuCategory, int64, error) {
	q := mt.c.MenuCategory.Query()
	if f.Name != "" {
		q = q.Where(menucategory.NameContains(f.Name))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if total <= 0 {
		return []*ent.MenuCategory{}, 0, nil
	}

	if f.Limit <= 0 {
		f.Limit = 10
	}

	list, err := q.Order(
		menucategory.BySort(sql.OrderAsc()),
		menucategory.ByID(sql.OrderAsc()),
	).Limit(f.Limit).Offset(f.Offset).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return list, int64(total), nil
}

// Update 更新分类（全量字段）。description 为空串时清空数据库中的描述。
func (mt *MenuType) Update(ctx context.Context, id uint64, name, description string, sort uint32) error {
	upd := mt.c.MenuCategory.UpdateOneID(id).SetName(name).SetSort(sort)
	if description != "" {
		upd = upd.SetDescription(description)
	} else {
		upd = upd.ClearDescription()
	}

	_, err := upd.Save(ctx)
	return err
}

// Delete 删除分类。
func (mt *MenuType) Delete(ctx context.Context, id uint64) error {
	return mt.c.MenuCategory.DeleteOneID(id).Exec(ctx)
}

// CountMenusByCategoryID 统计某分类下菜单数量。
func (mt *MenuType) CountMenusByCategoryID(ctx context.Context, categoryID uint64) (int, error) {
	return mt.c.Menu.Query().Where(entmenu.HasCategoryWith(menucategory.IDEQ(categoryID))).Count(ctx)
}
