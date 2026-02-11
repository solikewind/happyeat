// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/tablecategory"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 列出餐桌类别
func ListTableCategoriesHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ListTableCategoriesReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := tablecategory.NewListTableCategoriesLogic(r.Context(), svcCtx)
		resp, err := l.ListTableCategories(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
