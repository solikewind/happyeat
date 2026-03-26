// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menucategory

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/menucategory"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 更新菜单种类
func UpdateMenuCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UpdateMenuCategoryReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menucategory.NewUpdateMenuCategoryLogic(r.Context(), svcCtx)
		resp, err := l.UpdateMenuCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
