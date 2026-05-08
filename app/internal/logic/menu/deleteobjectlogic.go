// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteObjectLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除对象
func NewDeleteObjectLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteObjectLogic {
	return &DeleteObjectLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteObjectLogic) DeleteObject(req *types.DeleteObjectReq) (*types.DeleteObjectReply, error) {
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
	if _, err := l.svcCtx.Cos.Object.Delete(l.ctx, item.Key, &cos.ObjectDeleteOptions{}); err != nil {
		return nil, err
	}
	if err := l.svcCtx.Object.DeleteByID(l.ctx, req.Id); err != nil {
		return nil, err
	}
	return &types.DeleteObjectReply{}, nil
}
