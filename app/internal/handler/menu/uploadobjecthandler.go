// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"net/http"

	"github.com/solikewind/happyeat/app/internal/logic/menu"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 上传对象（用于菜单图片）
func UploadObjectHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		file, header, err := r.FormFile("file")
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}
		defer file.Close()

		l := menu.NewUploadObjectLogic(r.Context(), svcCtx)
		resp, err := l.UploadObject(file, header)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
