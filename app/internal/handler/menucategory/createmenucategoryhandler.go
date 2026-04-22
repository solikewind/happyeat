// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menucategory

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/handler/jsonbody"
	"github.com/solikewind/happyeat/app/internal/logic/menucategory"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 创建菜单种类
func CreateMenuCategoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CreateMenuCategoryReq
		if err := jsonbody.Decode(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := menucategory.NewCreateMenuCategoryLogic(r.Context(), svcCtx)
		resp, err := l.CreateMenuCategory(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
