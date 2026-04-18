package spec

import (
	"context"
	"errors"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/dal/model/ent"
	specmodel "github.com/solikewind/happyeat/dal/model/spec"
)

func ensureGlobalSpecItem(ctx context.Context, svcCtx *svc.ServiceContext, specType, specValue string, defaultPrice int64) (*ent.SpecItem, error) {
	groupName := strings.TrimSpace(specType)
	itemName := strings.TrimSpace(specValue)
	if groupName == "" || itemName == "" {
		return nil, nil
	}

	group, err := svcCtx.SpecGroup.GetByName(ctx, groupName)
	if err != nil {
		if !ent.IsNotFound(err) {
			return nil, err
		}
		group, err = svcCtx.SpecGroup.Create(ctx, specmodel.CreateSpecGroupInput{
			Name: groupName,
			Sort: 0,
		})
		if err != nil {
			if !ent.IsConstraintError(err) {
				return nil, err
			}
			group, err = svcCtx.SpecGroup.GetByName(ctx, groupName)
			if err != nil {
				return nil, err
			}
		}
	}

	item, err := svcCtx.SpecItem.GetByName(ctx, group.ID, itemName)
	if err == nil {
		return item, nil
	}
	if !ent.IsNotFound(err) {
		return nil, err
	}

	item, err = svcCtx.SpecItem.Create(ctx, specmodel.CreateSpecItemInput{
		SpecGroupID:  group.ID,
		Name:         itemName,
		DefaultPrice: defaultPrice,
	})
	if err != nil {
		if !ent.IsConstraintError(err) {
			return nil, err
		}
		return svcCtx.SpecItem.GetByName(ctx, group.ID, itemName)
	}

	return item, nil
}

func resolveCategorySpecSource(ctx context.Context, svcCtx *svc.ServiceContext, specItemID uint64, specType, specValue string, priceDelta int64) (uint64, string, string, error) {
	if specItemID > 0 {
		specItem, err := svcCtx.SpecItem.GetByID(ctx, specItemID)
		if err != nil {
			if ent.IsNotFound(err) {
				return 0, "", "", errors.New("规格项不存在")
			}
			return 0, "", "", err
		}
		specGroup, err := svcCtx.SpecGroup.GetByID(ctx, specItem.SpecGroupID)
		if err != nil {
			if ent.IsNotFound(err) {
				return 0, "", "", errors.New("规格组不存在")
			}
			return 0, "", "", err
		}
		return specItem.ID, specGroup.Name, specItem.Name, nil
	}

	specType = strings.TrimSpace(specType)
	specValue = strings.TrimSpace(specValue)
	if specType == "" {
		return 0, "", "", errors.New("规格类型不能为空")
	}
	if specValue == "" {
		return 0, "", "", errors.New("规格选项不能为空")
	}

	specItem, err := ensureGlobalSpecItem(ctx, svcCtx, specType, specValue, priceDelta)
	if err != nil {
		return 0, "", "", err
	}
	if specItem == nil {
		return 0, "", "", errors.New("规格项不能为空")
	}
	return specItem.ID, specType, specValue, nil
}
