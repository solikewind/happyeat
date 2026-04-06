// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package table

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/pkg/dberr"
	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	daltable "github.com/solikewind/happyeat/dal/model/table"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTableLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新餐桌
func NewUpdateTableLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTableLogic {
	return &UpdateTableLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTableLogic) UpdateTable(req *types.UpdateTableReq) (resp *types.UpdateTableReply, err error) {
	if req.Code == "" {
		return nil, errors.New("餐桌编号不能为空")
	}
	if req.Status == "" {
		return nil, errors.New("餐桌状态不能为空")
	}
	if req.Capacity <= 0 {
		return nil, errors.New("餐桌容量不能小于等于0")
	}
	if req.CategoryId == 0 {
		return nil, errors.New("请选择餐桌分类")
	}

	_, err = l.svcCtx.TableType.GetByID(l.ctx, req.CategoryId)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("餐桌分类不存在")
		}
		return nil, err
	}

	err = l.svcCtx.Table.Update(l.ctx, req.Id, daltable.UpdateTableInput{
		Code:       req.Code,
		Status:     req.Status,
		Capacity:   req.Capacity,
		CategoryID: req.CategoryId,
		QRCode:     req.QrCode,
	})
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("餐桌不存在")
		}
		if dberr.IsForeignKeyViolation(err) {
			return nil, errors.New("餐桌分类不存在或已被删除，请刷新分类列表后重选")
		}
		return nil, err
	}

	return &types.UpdateTableReply{}, nil
}
