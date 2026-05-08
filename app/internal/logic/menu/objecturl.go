package menu

import (
	"context"
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
