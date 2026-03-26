// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package spec

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/spec"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 删除规格项
func DeleteSpecItemHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DeleteSpecItemReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := spec.NewDeleteSpecItemLogic(r.Context(), svcCtx)
		resp, err := l.DeleteSpecItem(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
