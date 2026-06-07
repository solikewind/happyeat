// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"
	"errors"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	"github.com/solikewind/happyeat/dal/model/menu"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 更新菜单
func NewUpdateMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateMenuLogic {
	return &UpdateMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateMenuLogic) UpdateMenu(req *types.UpdateMenuReq) (*types.UpdateMenuReply, error) {
	if req.Name == "" {
		return nil, errors.New("菜单名称不能为空")
	}
	if req.Price < 0 {
		return nil, errors.New("价格不能为负")
	}

	existing, err := l.svcCtx.Menu.GetByID(l.ctx, req.Id)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("菜单不存在")
		}
		return nil, err
	}

	categoryID := req.CategoryId
	if categoryID == 0 {
		cat, errCat := existing.Edges.CategoryOrErr()
		if errCat != nil || cat == nil {
			return nil, errors.New("请指定菜单分类")
		}
		categoryID = cat.ID
	}

	_, err = l.svcCtx.MenuType.GetByID(l.ctx, categoryID)
	if err != nil {
		if ent.IsNotFound(err) {
			return nil, errors.New("分类不存在")
		}
		return nil, err
	}

	var specsPtr *[]menu.SpecInput
	if req.Specs != nil {
		specs, err := resolveMenuSpecs(l.ctx, l.svcCtx, categoryID, req.Specs)
		if err != nil {
			return nil, err
		}
		specsPtr = &specs
	}

	image := req.Image
	objectID := req.ObjectId

	var existingObjectID uint64
	if existing.ObjectID != nil {
		existingObjectID = *existing.ObjectID
	}

	finalObjectID := objectID
	if finalObjectID == 0 {
		finalObjectID = existingObjectID
	}

	if objectID > 0 {
		obj, err := l.svcCtx.Object.GetByID(l.ctx, objectID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errors.New("对象不存在")
			}
			return nil, err
		}
		if req.Image != "" && !sameObjectImageURL(req.Image, obj.URL) {
			return nil, errors.New("object_id 与 image 不一致")
		}
		image = obj.URL
	} else if finalObjectID > 0 && req.Image != "" {
		obj, err := l.svcCtx.Object.GetByID(l.ctx, finalObjectID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errors.New("关联对象不存在")
			}
			return nil, err
		}
		if !sameObjectImageURL(req.Image, obj.URL) {
			return nil, errors.New("已绑定 object_id 时，image 必须与对象 URL 一致，或改用 object_id 更新封面")
		}
		image = obj.URL
	} else if finalObjectID > 0 && req.Image == "" {
		obj, err := l.svcCtx.Object.GetByID(l.ctx, finalObjectID)
		if err != nil {
			if ent.IsNotFound(err) {
				return nil, errors.New("关联对象不存在")
			}
			return nil, err
		}
		image = obj.URL
	}

	err = l.svcCtx.Menu.Update(l.ctx, req.Id, menu.UpdateMenuInput{
		Name:        req.Name,
		Description: req.Description,
		Image:       image,
		ObjectID:    finalObjectID,
		Price:       req.Price,
		CategoryID:  categoryID,
		Sort:        req.Sort,
		Specs:       specsPtr,
	})
	if err != nil {
		return nil, err
	}

	return &types.UpdateMenuReply{}, nil
}
