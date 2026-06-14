// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package settlement

import (
	"context"
	"errors"
	"strings"

	"github.com/solikewind/happyeat/app/internal/svc"
	"github.com/solikewind/happyeat/app/internal/types"
	"github.com/solikewind/happyeat/common/consts/enum"
	dalsettlement "github.com/solikewind/happyeat/dal/model/settlement"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSettlementLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListSettlementLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSettlementLogic {
	return &ListSettlementLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListSettlementLogic) ListSettlement(req *types.ListSettlementReq) (*types.ListSettlementReply, error) {
	var status enum.SettlementStatus
	if s := strings.TrimSpace(req.Status); s != "" {
		parsed, err := parseSettlementStatus(s)
		if err != nil {
			return nil, err
		}
		status = parsed
	}

	pageSize := req.PageSize
	if pageSize <= 0 {
		pageSize = 10
	}
	current := req.Current
	if current <= 0 {
		current = 1
	}

	list, total, err := l.svcCtx.Settlement.List(l.ctx, dalsettlement.ListFilter{
		Status:       status,
		CustomerName: req.CustomerName,
		Offset:       int((current - 1) * pageSize),
		Limit:        int(pageSize),
	})
	if err != nil {
		return nil, err
	}

	out := make([]types.Settlement, 0, len(list))
	for _, item := range list {
		count, err := l.svcCtx.Settlement.CountOrders(l.ctx, item.ID)
		if err != nil {
			return nil, err
		}
		out = append(out, EntListItemToType(item, count))
	}

	return &types.ListSettlementReply{
		Settlements: out,
		Total:       total,
	}, nil
}

func parseSettlementStatus(s string) (enum.SettlementStatus, error) {
	switch strings.ToUpper(strings.TrimSpace(s)) {
	case string(enum.SettlementStatusUnsettled), "未结账":
		return enum.SettlementStatusUnsettled, nil
	case string(enum.SettlementStatusSettled), "已结账":
		return enum.SettlementStatusSettled, nil
	default:
		return "", errors.New("status 应为 UNSETTLED 或 SETTLED")
	}
}
