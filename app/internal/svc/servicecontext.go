// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"database/sql"

	"github.com/solikewind/happyeat/app/internal/config"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/menu"
	"github.com/solikewind/happyeat/dal/model/order"
	"github.com/solikewind/happyeat/dal/model/table"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB // 共享连接池，仅用于关闭
	Casbin *CasbinEnforcer // 权限：model 来自配置内联，policy 来自 DB casbin_rule 表

	Menu     *menu.Menu     // 菜单 data 层
	MenuType *menu.MenuType // 分类 data 层

	Table     *table.Table     // 餐桌 data 层
	TableType *table.TableType // 餐桌分类 data 层

	Order *order.Order // 订单 data 层
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	db, err := sql.Open("pgx", c.SqlConfig.DataSource)
	if err != nil {
		return nil, err
	}

	drv := entsql.OpenDB(dialect.Postgres, db)
	client := ent.NewClient(ent.Driver(drv))

	ce, err := NewCasbinEnforcer(c.Casbin.Model, c.SqlConfig.DataSource)
	if err != nil {
		return nil, err
	}

	return &ServiceContext{
		Config: c,
		DB:     db,
		Casbin: ce,

		Menu:      menu.NewMenu(client),
		MenuType:  menu.NewMenuType(client),
		Table:     table.NewTable(client),
		TableType: table.NewTableType(client),
		Order:     order.NewOrder(client),
	}, nil
}
