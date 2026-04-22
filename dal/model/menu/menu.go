package menu

import (
	"context"
	"strings"

	"entgo.io/ent/dialect/sql"
	"github.com/solikewind/happyeat/dal/model/ent"
	entmenu "github.com/solikewind/happyeat/dal/model/ent/menu"
	"github.com/solikewind/happyeat/dal/model/ent/menucategory"
	"github.com/solikewind/happyeat/dal/model/ent/menuspec"
)

type Menu struct {
	c *ent.Client
}

func NewMenu(c *ent.Client) *Menu {
	return &Menu{c: c}
}

type SpecInput struct {
	SpecItemID     uint64
	CategorySpecID uint64
	SpecType       string
	SpecValue      string
	PriceDelta     int64
	Sort           uint32
}

type CreateMenuInput struct {
	Name        string
	Description string
	Image       string
	Price       int64
	CategoryID  uint64
	Specs       []SpecInput
}

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
		createSpec := tx.MenuSpec.Create().
			SetMenuID(entMenu.ID).
			SetPriceDelta(spec.PriceDelta)
		if spec.SpecItemID > 0 {
			createSpec = createSpec.SetSpecItemID(spec.SpecItemID)
		}
		if spec.CategorySpecID > 0 {
			createSpec = createSpec.SetCategorySpecID(spec.CategorySpecID)
		}
		if spec.Sort > 0 {
			createSpec = createSpec.SetSort(spec.Sort)
		} else {
			createSpec = createSpec.SetSort(uint32(i))
		}
		_, err = createSpec.Save(ctx)
		if err != nil {
			return nil, err
		}
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return m.withMenuEdges(m.c.Menu.Query().Where(entmenu.IDEQ(entMenu.ID))).Only(ctx)
}

func (m *Menu) GetByID(ctx context.Context, id uint64) (*ent.Menu, error) {
	return m.withMenuEdges(m.c.Menu.Query().Where(entmenu.IDEQ(id))).Only(ctx)
}

func (m *Menu) Exist(ctx context.Context, id uint64) (bool, error) {
	return m.c.Menu.Query().Where(entmenu.IDEQ(id)).Exist(ctx)
}

type ListMenusFilter struct {
	Name         string
	CategoryName string
	CategoryID   uint64
	Offset       int
	Limit        int
}

const maxPinyinScanRows = 1000

func (m *Menu) List(ctx context.Context, f ListMenusFilter) ([]*ent.Menu, int64, error) {
	if f.Limit <= 0 {
		f.Limit = 10
	}
	if f.Offset < 0 {
		f.Offset = 0
	}

	newBaseQuery := func() *ent.MenuQuery {
		q := m.c.Menu.Query()
		if f.CategoryName != "" {
			q = q.Where(entmenu.HasCategoryWith(menucategory.NameEQ(f.CategoryName)))
		}
		if f.CategoryID > 0 {
			q = q.Where(entmenu.MenuCategoryIDEQ(f.CategoryID))
		}
		return q
	}

	name := strings.TrimSpace(f.Name)
	if name == "" {
		total, err := newBaseQuery().Count(ctx)
		if err != nil {
			return nil, 0, err
		}
		list, err := m.withMenuEdges(newBaseQuery()).
			Order(entmenu.ByID(sql.OrderDesc())).
			Offset(f.Offset).
			Limit(f.Limit).
			All(ctx)
		if err != nil {
			return nil, 0, err
		}
		return list, int64(total), nil
	}

	candidatesQuery := newBaseQuery().Where(entmenu.NameContainsFold(name))
	total, err := candidatesQuery.Count(ctx)
	if err != nil {
		return nil, 0, err
	}
	if total > 0 {
		list, err := m.withMenuEdges(newBaseQuery().
			Where(entmenu.NameContainsFold(name))).
			Order(entmenu.ByID(sql.OrderDesc())).
			Offset(f.Offset).
			Limit(f.Limit).
			All(ctx)
		if err != nil {
			return nil, 0, err
		}
		return list, int64(total), nil
	}

	// Fallback for pinyin keyword search; cap scan size to avoid expensive full-table scans.
	allList, err := m.withMenuEdges(newBaseQuery()).
		Order(entmenu.ByID(sql.OrderDesc())).
		Limit(maxPinyinScanRows).
		All(ctx)
	if err != nil {
		return nil, 0, err
	}

	matchedList := make([]*ent.Menu, 0, len(allList))
	for _, menu := range allList {
		if MatchPinyin(menu.Name, name) {
			matchedList = append(matchedList, menu)
		}
	}

	totalMatched := len(matchedList)
	if f.Offset >= totalMatched {
		return []*ent.Menu{}, int64(totalMatched), nil
	}
	end := f.Offset + f.Limit
	if end > totalMatched {
		end = totalMatched
	}
	return matchedList[f.Offset:end], int64(totalMatched), nil
}

type UpdateMenuInput struct {
	Name        string
	Description string
	Image       string
	Price       int64
	CategoryID  uint64
	Specs       []SpecInput
}

func (m *Menu) Update(ctx context.Context, id uint64, in UpdateMenuInput) error {
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
		createSpec := tx.MenuSpec.Create().
			SetMenuID(id).
			SetPriceDelta(spec.PriceDelta)
		if spec.SpecItemID > 0 {
			createSpec = createSpec.SetSpecItemID(spec.SpecItemID)
		}
		if spec.CategorySpecID > 0 {
			createSpec = createSpec.SetCategorySpecID(spec.CategorySpecID)
		}
		if spec.Sort > 0 {
			createSpec = createSpec.SetSort(spec.Sort)
		} else {
			createSpec = createSpec.SetSort(uint32(i))
		}
		_, err = createSpec.Save(ctx)
		if err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (m *Menu) Delete(ctx context.Context, id uint64) error {
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

func (m *Menu) withMenuEdges(q *ent.MenuQuery) *ent.MenuQuery {
	return q.WithCategory().WithMenuSpecs(func(sq *ent.MenuSpecQuery) {
		sq.Order(
			menuspec.BySort(sql.OrderAsc()),
			menuspec.ByID(sql.OrderAsc()),
		)
		sq.WithCategorySpec()
		sq.WithSpecItem(func(iq *ent.SpecItemQuery) {
			iq.WithSpecGroup()
		})
	})
}
