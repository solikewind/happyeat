package svc

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	pgxadapter "github.com/pckhoi/casbin-pgx-adapter/v3"
)

type CasbinEnforcer struct {
	Enforcer *casbin.Enforcer
}

// dbNameFromDataSource 从 postgres/postgresql URL 中解析出数据库名，与业务库一致，避免 adapter 使用默认的 "casbin" 库。
func dbNameFromDataSource(dataSource string) string {
	u, err := url.Parse(dataSource)
	if err != nil {
		return ""
	}
	// postgres://... 的 path 形如 "/happyeat"
	name := strings.TrimPrefix(u.Path, "/")
	if idx := strings.Index(name, "?"); idx >= 0 {
		name = name[:idx]
	}
	return name
}

// toPGXAdapterDSN 将 postgres URI 规范化为 key=value DSN，避免某些 adapter 对 URI 解析后回退 localhost。
func toPGXAdapterDSN(dataSource string) string {
	u, err := url.Parse(dataSource)
	if err != nil {
		return dataSource
	}
	if u.Scheme != "postgres" && u.Scheme != "postgresql" {
		return dataSource
	}

	host := u.Hostname()
	port := u.Port()
	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "5432"
	}

	user := ""
	pass := ""
	if u.User != nil {
		user = u.User.Username()
		pass, _ = u.User.Password()
	}
	db := strings.TrimPrefix(u.Path, "/")
	sslmode := u.Query().Get("sslmode")
	if sslmode == "" {
		sslmode = "disable"
	}

	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		host, port, user, pass, db, sslmode,
	)
}

// NewCasbinEnforcer 使用 yaml 中的 model 字符串与数据库（DataSource）中的 casbin_rule 表作为策略存储。
// 表建在 DataSource 指定的同一库（如 happyeat），表若不存在 adapter 会自动创建。
func NewCasbinEnforcer(modelText, dataSource string) (*CasbinEnforcer, error) {
	m, err := model.NewModelFromString(modelText)
	if err != nil {
		return nil, err
	}
	opts := []pgxadapter.Option{}
	if db := dbNameFromDataSource(dataSource); db != "" {
		opts = append(opts, pgxadapter.WithDatabase(db))
	}
	a, err := pgxadapter.NewAdapter(toPGXAdapterDSN(dataSource), opts...)
	if err != nil {
		return nil, err
	}
	e, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}
	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}
	return &CasbinEnforcer{Enforcer: e}, nil
}
