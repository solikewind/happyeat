package menu

import (
	"context"
	"strings"
	"time"

	"github.com/solikewind/happyeat/app/internal/svc"
)

const objectURLExpire = 10 * time.Minute

func signedURLOrRaw(ctx context.Context, svcCtx *svc.ServiceContext, key, rawURL string) string {
	if svcCtx == nil || svcCtx.Cos == nil || key == "" {
		return rawURL
	}
	signedURL, err := svcCtx.Cos.PresignedGetURL(ctx, key, objectURLExpire)
	if err != nil || signedURL == "" {
		return rawURL
	}
	return signedURL
}

// sameObjectImageURL 判断客户端回传的 image 是否指向同一对象（忽略 COS 签名等 query）。
func sameObjectImageURL(got, canonical string) bool {
	if got == "" || canonical == "" {
		return got == canonical
	}
	if got == canonical {
		return true
	}
	return stripURLQuery(got) == stripURLQuery(canonical)
}

func stripURLQuery(raw string) string {
	if i := strings.IndexByte(raw, '?'); i >= 0 {
		return raw[:i]
	}
	return raw
}
