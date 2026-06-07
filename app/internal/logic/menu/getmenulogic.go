// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package menu

import (
	"context"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/util/timeutil"
	"github.com/solikewind/happyeat/dal/model/ent"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMenuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 获取单个菜单
func NewGetMenuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMenuLogic {
	return &GetMenuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMenuLogic) GetMenu(req *types.GetMenuReq) (*types.GetMenuReply, error) {
	entMenu, err := l.svcCtx.Menu.GetByID(l.ctx, req.Id)
	if err != nil {
		return nil, err
	}

	return &types.GetMenuReply{Menu: entMenuToType(l.ctx, l.svcCtx, entMenu)}, nil
}

func entMenuToType(ctx context.Context, svcCtx *svc.ServiceContext, e *ent.Menu) types.Menu {
	out := types.Menu{
		Id:        uint64(e.ID),
		Name:      e.Name,
		Price:     e.Price,
		Sort:      e.Sort,
		CreatedAt: timeutil.TimeToString(e.CreatedAt),
		UpdatedAt: timeutil.TimeToString(e.UpdatedAt),
	}

	if e.Description != nil {
		out.Description = *e.Description
	}
	if e.ObjectID != nil && *e.ObjectID > 0 {
		out.ObjectId = *e.ObjectID
	}
	coverObj, _ := e.Edges.CoverObjectOrErr()
	imageFromField := ""
	if e.Image != nil {
		imageFromField = *e.Image
	}
	if coverObj != nil {
		preferURL := signedURLOrRaw(ctx, svcCtx, coverObj.Key, coverObj.URL)
		if imageFromField != "" && imageFromField != coverObj.URL {
			// 数据不一致时仍以对象 URL 为准（object_id 为权威来源）
			out.Image = preferURL
		} else if imageFromField == "" {
			out.Image = preferURL
		} else {
			out.Image = signedURLOrRaw(ctx, svcCtx, coverObj.Key, imageFromField)
		}
	} else if imageFromField != "" {
		out.Image = imageFromField
	}

	cat, _ := e.Edges.CategoryOrErr()
	if cat != nil {
		out.CategoryId = uint64(cat.ID)
	}

	specs, _ := e.Edges.MenuSpecsOrErr()
	for _, s := range specs {
		spec := types.MenuSpec{
			PriceDelta: s.PriceDelta,
			Sort:       s.Sort,
		}
		if s.SpecItemID != nil {
			spec.SpecItemId = *s.SpecItemID
		}
		if s.CategorySpecID != nil {
			spec.CategorySpecId = *s.CategorySpecID
			spec.Source = "category"
		}
		if categorySpec, err := s.Edges.CategorySpecOrErr(); err == nil && categorySpec != nil {
			spec.SpecType = categorySpec.SpecType
			spec.SpecValue = categorySpec.SpecValue
		}
		if specItem, err := s.Edges.SpecItemOrErr(); err == nil && specItem != nil {
			if spec.Source == "" {
				spec.Source = "library"
			}
			if spec.SpecValue == "" {
				spec.SpecValue = specItem.Name
			}
			if spec.SpecType == "" {
				if group, groupErr := specItem.Edges.SpecGroupOrErr(); groupErr == nil && group != nil {
					spec.SpecType = group.Name
				}
			}
		}
		out.Specs = append(out.Specs, spec)
	}

	return out
}
