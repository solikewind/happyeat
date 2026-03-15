// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package tablecategory

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTableCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新餐桌类别
func NewUpdateTableCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTableCategoryLogic {
	return &UpdateTableCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateTableCategoryLogic) UpdateTableCategory(req *types.UpdateTableCategoryReq) (resp *types.UpdateTableCategoryReply, err error) {
	c := req.TableCategory
	if c.Name == "" {
		return nil, errors.New("分类名称不能为空")
	}

	err = l.svcCtx.TableType.Update(l.ctx, int(req.Id), c.Name, c.Description)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("餐桌分类不存在")
		}
		return nil, err
	}

	return &types.UpdateTableCategoryReply{}, nil
}
