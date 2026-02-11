// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"database/sql"

	"github.com/solikewind/happyeat/app/internal/config"
	"github.com/solikewind/happyeat/dal/model/menu"
	entmenu "github.com/solikewind/happyeat/dal/model/menu/ent"
	"github.com/solikewind/happyeat/dal/model/table"
	enttable "github.com/solikewind/happyeat/dal/model/table/ent"
	"entgo.io/ent/dialect"
	entsql "entgo.io/ent/dialect/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type ServiceContext struct {
	Config config.Config
	DB     *sql.DB // 共享连接池，仅用于关闭

	Menu     *menu.Menu     // 菜单 data 层
	MenuType *menu.MenuType // 分类 data 层

	Table     *table.Table     // 餐桌 data 层
	TableType *table.TableType // 餐桌分类 data 层
}

func NewServiceContext(c config.Config) (*ServiceContext, error) {
	db, err := sql.Open("pgx", c.SqlConfig.DataSource)
	if err != nil {
		return nil, err
	}

	drvMenu := entsql.OpenDB(dialect.Postgres, db)
	drvTable := entsql.OpenDB(dialect.Postgres, db)
	clientMenu := entmenu.NewClient(entmenu.Driver(drvMenu))
	clientTable := enttable.NewClient(enttable.Driver(drvTable))

	return &ServiceContext{
		Config: c,
		DB:     db,

		Menu:      menu.NewMenu(clientMenu),
		MenuType:  menu.NewMenuType(clientMenu),
		Table:     table.NewTable(clientTable),
		TableType: table.NewTableType(clientTable),
	}, nil
}
