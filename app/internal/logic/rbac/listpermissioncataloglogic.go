// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package rbac

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/pkg/casbinrules"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPermissionCatalogLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 权限目录（权限码及对应 HTTP 接口，供前端勾选）
func NewListPermissionCatalogLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPermissionCatalogLogic {
	return &ListPermissionCatalogLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListPermissionCatalogLogic) ListPermissionCatalog() (resp *types.ListPermissionCatalogReply, err error) {
	catalog := casbinrules.ListCatalog()
	items := make([]types.PermissionCatalogItem, 0, len(catalog))
	for _, entry := range catalog {
		eps := make([]types.PermissionEndpoint, 0, len(entry.Endpoints))
		for _, ep := range entry.Endpoints {
			eps = append(eps, types.PermissionEndpoint{Obj: ep.Obj, Act: ep.Act})
		}
		items = append(items, types.PermissionCatalogItem{
			Code:        entry.Code,
			Description: entry.Description,
			Endpoints:   eps,
		})
	}
	return &types.ListPermissionCatalogReply{Permissions: items}, nil
}
