// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package MenuCategory

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

const (
	errCategoryHasMenus  = "该分类下仍有菜单，无法删除"
	errInvalidCategoryID = "无效的分类ID"
)

var (
	ErrCategoryHasMenus  = errors.New(errCategoryHasMenus)
	ErrInvalidCategoryID = errors.New(errInvalidCategoryID)
)

type DeleteMenuCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 删除菜单种类
func NewDeleteMenuCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteMenuCategoryLogic {
	return &DeleteMenuCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteMenuCategoryLogic) DeleteMenuCategory(req *types.DeleteMenuCategoryReq) (*types.DeleteMenuCategoryReply, error) {
	// 参数验证
	if req.Id == 0 {
		return nil, ErrInvalidCategoryID
	}

	l.Infof("开始删除菜单分类，ID: %d", req.Id)

	count, err := l.svcCtx.MenuType.CountMenusByCategoryID(l.ctx, int(req.Id))
	if err != nil {
		l.Errorf("查询分类菜单数量失败，ID: %d, 错误: %v", req.Id, err)
		return nil, err
	}

	if count > 0 {
		l.Errorf("无法删除菜单分类，该分类下仍有菜单，ID: %d, 菜单数: %d", req.Id, count)
		return nil, ErrCategoryHasMenus
	}

	err = l.svcCtx.MenuType.Delete(l.ctx, int(req.Id))
	if err != nil {
		l.Errorf("删除菜单分类失败，ID: %d, 错误: %v", req.Id, err)
		return nil, err
	}

	l.Infof("成功删除菜单分类，ID: %d", req.Id)
	return &types.DeleteMenuCategoryReply{}, nil
}
