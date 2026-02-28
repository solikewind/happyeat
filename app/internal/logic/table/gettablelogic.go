// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个餐桌
func NewGetTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetTableLogic {
	return &GetTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetTableLogic) GetTable(req *types.GetTableReq) (resp *types.GetTableReply, err error) {
	e, err := l.svcCtx.Table.GetByID(l.ctx, int(req.Id))
	if err != nil {
		l.Errorf("GetTable GetByID err: %v", err)
		return nil, err
	}
	categoryID := uint64(0)
	if e.Edges.Category != nil {
		categoryID = uint64(e.Edges.Category.ID)
	}
	qrCode := ""
	if e.QrCode != nil {
		qrCode = *e.QrCode
	}
	return &types.GetTableReply{
		Table: types.Table{
			Id:         uint64(e.ID),
			Code:       e.Code,
			Status:     e.Status,
			Capacity:   e.Capacity,
			CategoryId: categoryID,
			QrCode:     qrCode,
			CreateAt:   e.CreatedAt.Unix(),
			UpdateAt:   e.UpdatedAt.Unix(),
		},
	}, nil
}
