package menu

import (
	"context"
	"errors"
	"time"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetObjectURLLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取对象临时访问地址（私有桶）
func NewGetObjectURLLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetObjectURLLogic {
	return &GetObjectURLLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetObjectURLLogic) GetObjectURL(req *types.GetObjectURLReq) (*types.GetObjectURLReply, error) {
	item, err := l.svcCtx.Object.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("对象不存在")
		}
		return nil, err
	}
	if l.svcCtx.Cos == nil {
		return nil, errors.New("cos 未配置")
	}
	expiredIn := objectURLExpire
	signedURL, err := l.svcCtx.Cos.PresignedGetURL(l.ctx, item.Key, expiredIn)
	if err != nil {
		return nil, err
	}
	return &types.GetObjectURLReply{
		Url:       signedURL,
		ExpiredAt: timeutil.TimeToString(time.Now().UTC().Add(expiredIn)),
	}, nil
}
