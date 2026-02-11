// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package homestayOrder

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/order"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/rpc/travel"
	"github.com/wwwzy/ZeroMicroServices/pkg/ctxdata"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateHomestayOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 创建民宿订单
func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateHomestayOrderLogic) CreateHomestayOrder(req *types.CreateHomestayOrderReq) (*types.CreateHomestayOrderResp, error) {
	//查询民宿详情
	homestayResp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		Id: req.HomestayId,
	})
	if err != nil {
		return nil, err
	}
	if homestayResp.Homestay == nil || homestayResp.Homestay.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("homestay no exists"), "CreateHomestayOrder homestay no exists id : %d", req.HomestayId)
	}

	//获取用户会话id并创建订单
	userId := ctxdata.GetUidFromCtx(l.ctx)
	resp, err := l.svcCtx.OrderRpc.CreateHomestayOrder(l.ctx, &order.CreateHomestayOrderReq{
		HomestayId:    req.HomestayId,
		IsFood:        req.IsFood,
		LiveStartTime: req.LiveStartTime,
		LiveEndTime:   req.LiveEndTime,
		UserId:        userId,
		LivePeopleNum: req.LivePeopleNum,
		Remark:        req.Remark,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("create homestay order fail"), "create homestay order rpc CreateHomestayOrder fail req: %+v , err : %v ", req, err)
	}

	return &types.CreateHomestayOrderResp{
		OrderSn: resp.Sn,
	}, nil
}
