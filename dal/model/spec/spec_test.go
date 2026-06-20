package spec

import (
	"context"
	"database/sql"
	"testing"

	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/ent/enttest"
	_ "modernc.org/sqlite"
)

func newTestClient(t *testing.T) *ent.Client {
	t.Helper()

	db, err := sql.Open("sqlite", "file:spec?mode=memory&cache=shared&_fk=1")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		t.Fatalf("enable foreign keys: %v", err)
	}

	drv := entsql.OpenDB(dialect.SQLite, db)
	return enttest.NewClient(t, enttest.WithOptions(ent.Driver(drv)))
}

func TestCategorySpecUpdateSyncsLinkedMenuSpecPriceDelta(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	category := client.MenuCategory.Create().
		SetName("本店特色").
		SaveX(ctx)
	categorySpec := client.CategorySpec.Create().
		SetMenuCategoryID(category.ID).
		SetSpecType("份量").
		SetSpecValue("大份").
		SetPriceDelta(1).
		SaveX(ctx)
	menu := client.Menu.Create().
		SetMenuCategoryID(category.ID).
		SetName("炒鸡").
		SetPrice(68).
		SaveX(ctx)
	menuSpec := client.MenuSpec.Create().
		SetMenuID(menu.ID).
		SetCategorySpecID(categorySpec.ID).
		SetPriceDelta(1).
		SaveX(ctx)

	err := NewCategorySpec(client).Update(ctx, categorySpec.ID, CreateCategorySpecInput{
		CategoryID: category.ID,
		SpecType:   "份量",
		SpecValue:  "大份",
		PriceDelta: 2,
		Sort:       categorySpec.Sort,
		SpecItemID: 0,
	})
	if err != nil {
		t.Fatalf("update category spec: %v", err)
	}

	updated := client.MenuSpec.GetX(ctx, menuSpec.ID)
	if updated.PriceDelta != 2 {
		t.Fatalf("menu spec price_delta = %d, want 2", updated.PriceDelta)
	}
}

func TestSpecItemUpdateSyncsLinkedSpecPrices(t *testing.T) {
	ctx := context.Background()
	client := newTestClient(t)
	defer client.Close()

	group := client.SpecGroup.Create().
		SetName("份量").
		SaveX(ctx)
	item := client.SpecItem.Create().
		SetSpecGroupID(group.ID).
		SetName("大份").
		SetDefaultPrice(1).
		SaveX(ctx)
	category := client.MenuCategory.Create().
		SetName("本店特色").
		SaveX(ctx)
	categorySpec := client.CategorySpec.Create().
		SetMenuCategoryID(category.ID).
		SetSpecItemID(item.ID).
		SetSpecType("份量").
		SetSpecValue("大份").
		SetPriceDelta(1).
		SaveX(ctx)
	menu := client.Menu.Create().
		SetMenuCategoryID(category.ID).
		SetName("炒鸡").
		SetPrice(68).
		SaveX(ctx)
	directMenuSpec := client.MenuSpec.Create().
		SetMenuID(menu.ID).
		SetSpecItemID(item.ID).
		SetPriceDelta(1).
		SaveX(ctx)
	categoryMenuSpec := client.MenuSpec.Create().
		SetMenuID(menu.ID).
		SetCategorySpecID(categorySpec.ID).
		SetPriceDelta(1).
		SaveX(ctx)

	err := NewSpecItem(client).Update(ctx, item.ID, CreateSpecItemInput{
		SpecGroupID:  group.ID,
		Name:         "大份",
		DefaultPrice: 2,
		Sort:         item.Sort,
	})
	if err != nil {
		t.Fatalf("update spec item: %v", err)
	}

	updatedCategorySpec := client.CategorySpec.GetX(ctx, categorySpec.ID)
	if updatedCategorySpec.PriceDelta != 2 {
		t.Fatalf("category spec price_delta = %d, want 2", updatedCategorySpec.PriceDelta)
	}
	updatedDirectMenuSpec := client.MenuSpec.GetX(ctx, directMenuSpec.ID)
	if updatedDirectMenuSpec.PriceDelta != 2 {
		t.Fatalf("direct menu spec price_delta = %d, want 2", updatedDirectMenuSpec.PriceDelta)
	}
	updatedCategoryMenuSpec := client.MenuSpec.GetX(ctx, categoryMenuSpec.ID)
	if updatedCategoryMenuSpec.PriceDelta != 2 {
		t.Fatalf("category menu spec price_delta = %d, want 2", updatedCategoryMenuSpec.PriceDelta)
	}
}
