package spec

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/ent/categoryspec"
	"github.com/solikewind/happyeat/dal/model/ent/menuspec"
	"github.com/solikewind/happyeat/dal/model/ent/specgroup"
	"github.com/solikewind/happyeat/dal/model/ent/specitem"
)

type CategorySpec struct {
	c *ent.Client
}

func NewCategorySpec(c *ent.Client) *CategorySpec {
	return &CategorySpec{c: c}
}

type CreateCategorySpecInput struct {
	CategoryID uint64
	SpecItemID uint64
	SpecType   string
	SpecValue  string
	PriceDelta int64
	Sort       uint32
}

type ListCategorySpecsFilter struct {
	CategoryID uint64
	SpecType   string
	Offset     int
	Limit      int
}

func (r *CategorySpec) Create(ctx context.Context, in CreateCategorySpecInput) (*ent.CategorySpec, error) {
	create := r.c.CategorySpec.Create().
		SetMenuCategoryID(in.CategoryID).
		SetSpecType(in.SpecType).
		SetSpecValue(in.SpecValue).
		SetPriceDelta(in.PriceDelta).
		SetSort(in.Sort)
	if in.SpecItemID > 0 {
		create = create.SetSpecItemID(in.SpecItemID)
	}
	return create.Save(ctx)
}

func (r *CategorySpec) GetByID(ctx context.Context, id uint64) (*ent.CategorySpec, error) {
	return r.c.CategorySpec.Query().
		Where(categoryspec.IDEQ(id)).
		WithSpecItem(func(iq *ent.SpecItemQuery) {
			iq.WithSpecGroup()
		}).
		Only(ctx)
}

func (r *CategorySpec) Exist(ctx context.Context, id uint64) (bool, error) {
	return r.c.CategorySpec.Query().Where(categoryspec.IDEQ(id)).Exist(ctx)
}

func (r *CategorySpec) ExistByValue(ctx context.Context, categoryID uint64, specType, specValue string, excludeID uint64) (bool, error) {
	query := r.c.CategorySpec.Query().Where(
		categoryspec.MenuCategoryIDEQ(categoryID),
		categoryspec.SpecTypeEQ(strings.TrimSpace(specType)),
		categoryspec.SpecValueEQ(strings.TrimSpace(specValue)),
	)
	if excludeID > 0 {
		query = query.Where(categoryspec.IDNEQ(excludeID))
	}
	return query.Exist(ctx)
}

func (r *CategorySpec) GetByValue(ctx context.Context, categoryID uint64, specType, specValue string) (*ent.CategorySpec, error) {
	return r.c.CategorySpec.Query().Where(
		categoryspec.MenuCategoryIDEQ(categoryID),
		categoryspec.SpecTypeEQ(strings.TrimSpace(specType)),
		categoryspec.SpecValueEQ(strings.TrimSpace(specValue)),
	).Only(ctx)
}

func (r *CategorySpec) List(ctx context.Context, f ListCategorySpecsFilter) ([]*ent.CategorySpec, int64, error) {
	q := r.c.CategorySpec.Query()
	if f.CategoryID > 0 {
		q = q.Where(categoryspec.MenuCategoryIDEQ(f.CategoryID))
	}
	if keyword := strings.TrimSpace(f.SpecType); keyword != "" {
		q = q.Where(categoryspec.SpecTypeContainsFold(keyword))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*ent.CategorySpec{}, 0, nil
	}

	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	list, err := q.WithSpecItem(func(iq *ent.SpecItemQuery) {
		iq.WithSpecGroup()
	}).Order(
		categoryspec.BySort(sql.OrderAsc()),
		categoryspec.ByID(sql.OrderAsc()),
	).Offset(f.Offset).Limit(f.Limit).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return list, int64(total), nil
}

func (r *CategorySpec) Update(ctx context.Context, id uint64, in CreateCategorySpecInput) error {
	upd := r.c.CategorySpec.UpdateOneID(id).
		SetMenuCategoryID(in.CategoryID).
		SetSpecType(in.SpecType).
		SetSpecValue(in.SpecValue).
		SetPriceDelta(in.PriceDelta).
		SetSort(in.Sort)
	if in.SpecItemID > 0 {
		upd = upd.SetSpecItemID(in.SpecItemID)
	} else {
		upd = upd.ClearSpecItem()
	}
	_, err := upd.Save(ctx)
	return err
}

func (r *CategorySpec) Delete(ctx context.Context, id uint64) error {
	return r.c.CategorySpec.DeleteOneID(id).Exec(ctx)
}

// CountMenuSpecsByCategory 统计有多少菜品规格引用了该分类规格模板
func (r *CategorySpec) CountMenuSpecsByCategory(ctx context.Context, categorySpecID uint64) (int, error) {
	return r.c.MenuSpec.Query().Where(menuspec.CategorySpecIDEQ(categorySpecID)).Count(ctx)
}

type SpecGroup struct {
	c *ent.Client
}

func NewSpecGroup(c *ent.Client) *SpecGroup {
	return &SpecGroup{c: c}
}

type CreateSpecGroupInput struct {
	Name string
	Sort uint32
}

type ListSpecGroupsFilter struct {
	Name   string
	Offset int
	Limit  int
}

func (r *SpecGroup) Create(ctx context.Context, in CreateSpecGroupInput) (*ent.SpecGroup, error) {
	return r.c.SpecGroup.Create().
		SetName(in.Name).
		SetSort(in.Sort).
		Save(ctx)
}

func (r *SpecGroup) GetByID(ctx context.Context, id uint64) (*ent.SpecGroup, error) {
	return r.c.SpecGroup.Get(ctx, id)
}

func (r *SpecGroup) GetByName(ctx context.Context, name string) (*ent.SpecGroup, error) {
	return r.c.SpecGroup.Query().Where(specgroup.NameEQ(strings.TrimSpace(name))).Only(ctx)
}

func (r *SpecGroup) Exist(ctx context.Context, id uint64) (bool, error) {
	return r.c.SpecGroup.Query().Where(specgroup.IDEQ(id)).Exist(ctx)
}

func (r *SpecGroup) ExistByName(ctx context.Context, name string, excludeID uint64) (bool, error) {
	query := r.c.SpecGroup.Query().Where(specgroup.NameEQ(strings.TrimSpace(name)))
	if excludeID > 0 {
		query = query.Where(specgroup.IDNEQ(excludeID))
	}
	return query.Exist(ctx)
}

func (r *SpecGroup) List(ctx context.Context, f ListSpecGroupsFilter) ([]*ent.SpecGroup, int64, error) {
	q := r.c.SpecGroup.Query()
	if keyword := strings.TrimSpace(f.Name); keyword != "" {
		q = q.Where(specgroup.NameContainsFold(keyword))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*ent.SpecGroup{}, 0, nil
	}

	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	list, err := q.Order(
		specgroup.BySort(sql.OrderAsc()),
		specgroup.ByID(sql.OrderAsc()),
	).Offset(f.Offset).Limit(f.Limit).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return list, int64(total), nil
}

func (r *SpecGroup) Update(ctx context.Context, id uint64, in CreateSpecGroupInput) error {
	_, err := r.c.SpecGroup.UpdateOneID(id).
		SetName(in.Name).
		SetSort(in.Sort).
		Save(ctx)
	return err
}

func (r *SpecGroup) Delete(ctx context.Context, id uint64) error {
	return r.c.SpecGroup.DeleteOneID(id).Exec(ctx)
}

func (r *SpecGroup) CountItemsByGroupID(ctx context.Context, groupID uint64) (int, error) {
	return r.c.SpecItem.Query().Where(specitem.SpecGroupIDEQ(groupID)).Count(ctx)
}

type SpecItem struct {
	c *ent.Client
}

func NewSpecItem(c *ent.Client) *SpecItem {
	return &SpecItem{c: c}
}

type CreateSpecItemInput struct {
	SpecGroupID  uint64
	Name         string
	DefaultPrice int64
	Sort         uint32
}

type ListSpecItemsFilter struct {
	SpecGroupID uint64
	Name        string
	Offset      int
	Limit       int
}

func (r *SpecItem) Create(ctx context.Context, in CreateSpecItemInput) (*ent.SpecItem, error) {
	return r.c.SpecItem.Create().
		SetSpecGroupID(in.SpecGroupID).
		SetName(in.Name).
		SetDefaultPrice(in.DefaultPrice).
		SetSort(in.Sort).
		Save(ctx)
}

func (r *SpecItem) GetByID(ctx context.Context, id uint64) (*ent.SpecItem, error) {
	return r.c.SpecItem.Get(ctx, id)
}

func (r *SpecItem) GetByName(ctx context.Context, groupID uint64, name string) (*ent.SpecItem, error) {
	return r.c.SpecItem.Query().Where(
		specitem.SpecGroupIDEQ(groupID),
		specitem.NameEQ(strings.TrimSpace(name)),
	).Only(ctx)
}

func (r *SpecItem) Exist(ctx context.Context, id uint64) (bool, error) {
	return r.c.SpecItem.Query().Where(specitem.IDEQ(id)).Exist(ctx)
}

func (r *SpecItem) ExistByName(ctx context.Context, groupID uint64, name string, excludeID uint64) (bool, error) {
	query := r.c.SpecItem.Query().Where(
		specitem.SpecGroupIDEQ(groupID),
		specitem.NameEQ(strings.TrimSpace(name)),
	)
	if excludeID > 0 {
		query = query.Where(specitem.IDNEQ(excludeID))
	}
	return query.Exist(ctx)
}

func (r *SpecItem) List(ctx context.Context, f ListSpecItemsFilter) ([]*ent.SpecItem, int64, error) {
	q := r.c.SpecItem.Query()
	if f.SpecGroupID > 0 {
		q = q.Where(specitem.SpecGroupIDEQ(f.SpecGroupID))
	}
	if keyword := strings.TrimSpace(f.Name); keyword != "" {
		q = q.Where(specitem.NameContainsFold(keyword))
	}

	total, err := q.Clone().Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if total == 0 {
		return []*ent.SpecItem{}, 0, nil
	}

	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	list, err := q.Order(
		specitem.BySort(sql.OrderAsc()),
		specitem.BySpecGroupID(sql.OrderAsc()),
		specitem.ByID(sql.OrderAsc()),
	).Offset(f.Offset).Limit(f.Limit).All(ctx)
	if err != nil {
		return nil, 0, err
	}

	return list, int64(total), nil
}

func (r *SpecItem) Update(ctx context.Context, id uint64, in CreateSpecItemInput) error {
	_, err := r.c.SpecItem.UpdateOneID(id).
		SetSpecGroupID(in.SpecGroupID).
		SetName(in.Name).
		SetDefaultPrice(in.DefaultPrice).
		SetSort(in.Sort).
		Save(ctx)
	return err
}

func (r *SpecItem) Delete(ctx context.Context, id uint64) error {
	return r.c.SpecItem.DeleteOneID(id).Exec(ctx)
}

// CountCategorySpecsByItem 统计有多少分类规格模板引用了该规格项
func (r *SpecItem) CountCategorySpecsByItem(ctx context.Context, specItemID uint64) (int, error) {
	return r.c.CategorySpec.Query().Where(categoryspec.SpecItemIDEQ(specItemID)).Count(ctx)
}

// CountMenuSpecsByItem 统计有多少菜品规格行直接引用了该规格项（不经过分类规格）
func (r *SpecItem) CountMenuSpecsByItem(ctx context.Context, specItemID uint64) (int, error) {
	return r.c.MenuSpec.Query().Where(menuspec.SpecItemIDEQ(specItemID)).Count(ctx)
}
