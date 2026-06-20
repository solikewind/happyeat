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
