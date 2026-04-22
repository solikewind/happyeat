package menu

import (
	"context"
	"errors"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/dal/model/ent"
	menumodel "github.com/solikewind/happyeat/dal/model/menu"
	specmodel "github.com/solikewind/happyeat/dal/model/spec"
)

// resolveMenuSpecs 解析请求规格列表，处理分类规格匹配和全局规格项的"找或建"，并去重。
//
// 去重策略（三路，命中任意一路即跳过）：
//   - 相同 category_spec_id
//   - 相同 spec_item_id（无 category_spec_id 时生效）
//   - 相同 spec_type::spec_value（兜底）
//
// source 语义：
//   - "category"  : 已携带 category_spec_id，跳过分类规格自动匹配
//   - "library"   : 已携带 spec_item_id，跳过分类规格自动匹配（避免将全局规格悄悄关联到分类规格）
//   - "custom"    : 无 ID，跳过分类规格匹配，但仍通过 ensureGlobalSpecItem 写入全局库
//   - ""（空）    : 自动尝试匹配分类规格，匹配不到则写全局库
func resolveMenuSpecs(ctx context.Context, svcCtx *svc.ServiceContext, categoryID uint64, reqSpecs []types.MenuSpec) ([]menumodel.SpecInput, error) {
	seenKey := map[string]bool{}     // spec_type::spec_value
	seenCatSpec := map[uint64]bool{} // category_spec_id
	seenSpecItem := map[uint64]bool{} // spec_item_id（仅在无 category_spec_id 时登记）

	specs := make([]menumodel.SpecInput, 0, len(reqSpecs))

	for _, s := range reqSpecs {
		si := menumodel.SpecInput{
			SpecItemID:     s.SpecItemId,
			CategorySpecID: s.CategorySpecId,
			SpecType:       strings.TrimSpace(s.SpecType),
			SpecValue:      strings.TrimSpace(s.SpecValue),
			PriceDelta:     s.PriceDelta,
			Sort:           s.Sort,
		}

		// 仅当未携带任何 ID 且 source 不是 library/custom 时，才尝试匹配分类规格模板
		if si.CategorySpecID == 0 && si.SpecItemID == 0 &&
			s.Source != "custom" && s.Source != "library" &&
			categoryID > 0 && si.SpecType != "" && si.SpecValue != "" {
			if cs, findErr := svcCtx.CategorySpec.GetByValue(ctx, categoryID, si.SpecType, si.SpecValue); findErr == nil && cs != nil {
				si.CategorySpecID = cs.ID
				si.SpecType = cs.SpecType
				si.SpecValue = cs.SpecValue
				if si.PriceDelta == 0 {
					si.PriceDelta = cs.PriceDelta
				}
			}
		}

		// 仍无任何 ID 时，按名称确保全局规格项（找或建）
		if si.CategorySpecID == 0 && si.SpecItemID == 0 && si.SpecType != "" && si.SpecValue != "" {
			item, ensureErr := ensureGlobalSpecItem(ctx, svcCtx, si.SpecType, si.SpecValue, si.PriceDelta)
			if ensureErr != nil {
				return nil, ensureErr
			}
			if item != nil {
				si.SpecItemID = item.ID
			}
		}

		if si.CategorySpecID == 0 && si.SpecItemID == 0 {
			return nil, errors.New("菜单规格必须引用分类规格或全局规格项")
		}

		// —— 去重 ——
		key := si.SpecType + "::" + si.SpecValue
		isDup := (key != "::" && seenKey[key]) ||
			(si.CategorySpecID > 0 && seenCatSpec[si.CategorySpecID]) ||
			(si.SpecItemID > 0 && si.CategorySpecID == 0 && seenSpecItem[si.SpecItemID])
		if isDup {
			continue
		}
		if key != "::" {
			seenKey[key] = true
		}
		if si.CategorySpecID > 0 {
			seenCatSpec[si.CategorySpecID] = true
		}
		if si.SpecItemID > 0 && si.CategorySpecID == 0 {
			seenSpecItem[si.SpecItemID] = true
		}

		specs = append(specs, si)
	}

	return specs, nil
}

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
