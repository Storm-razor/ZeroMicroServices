package logic

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/pb"
	"github.com/wwwzy/ZeroMicroServices/app/order/model"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/rpc/travel"
	"github.com/wwwzy/ZeroMicroServices/pkg/tool"
	"github.com/wwwzy/ZeroMicroServices/pkg/uniqueid"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

const CloseOrderTimeMinutes = 30 //订单关闭时间

type CreateHomestayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateHomestayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateHomestayOrderLogic {
	return &CreateHomestayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 民宿下订单服务
func (l *CreateHomestayOrderLogic) CreateHomestayOrder(in *pb.CreateHomestayOrderReq) (*pb.CreateHomestayOrderResp, error) {

	if in.LiveEndTime <= in.LiveStartTime {
		return nil, errors.Wrapf(xerr.NewErrMsg("Stay at least one night"), "Place an order at a B&B. The end time of your stay must be greater than the start time. in : %+v", in)
	}

	resp, err := l.svcCtx.TravelRpc.HomestayDetail(l.ctx, &travel.HomestayDetailReq{
		Id: in.HomestayId,
	})
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to query the record"), "Failed to query the record  rpc HomestayDetail fail , homestayId : %d , err : %v", in.HomestayId, err)
	}
	if resp.Homestay == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("This record does not exist"), "This record does not exist , homestayId : %d ", in.HomestayId)
	}

	//从民宿的轮播图列表中提取第一张图作为订单封面图
	var cover string
	if len(resp.Homestay.Banner) > 0 {
		cover = strings.Split(resp.Homestay.Banner, ",")[0]
	}

	//创建订单信息
	order := new(model.HomestayOrder)
	order.Sn = uniqueid.GenSn(uniqueid.SN_PREFIX_HOMESTAY_ORDER)
	order.UserId = in.UserId
	order.HomestayId = in.HomestayId
	order.Title = resp.Homestay.Title
	order.SubTitle = resp.Homestay.SubTitle
	order.Cover = cover
	order.Info = resp.Homestay.Info
	order.PeopleNum = resp.Homestay.PeopleNum
	order.RowType = resp.Homestay.RowType
	order.HomestayPrice = resp.Homestay.HomestayPrice
	order.MarketHomestayPrice = resp.Homestay.MarketHomestayPrice
	order.HomestayBusinessId = resp.Homestay.HomestayBusinessId
	order.HomestayUserId = resp.Homestay.UserId
	order.LivePeopleNum = in.LivePeopleNum
	order.TradeState = model.HomestayOrderTradeStateWaitPay
	order.TradeCode = tool.Krand(8, tool.KC_RAND_KIND_ALL)
	order.Remark = in.Remark
	order.FoodInfo = resp.Homestay.FoodInfo
	order.FoodPrice = resp.Homestay.FoodPrice
	order.LiveStartDate = time.Unix(in.LiveStartTime, 0)
	order.LiveEndDate = time.Unix(in.LiveEndTime, 0)

	liveDays := int64(order.LiveEndDate.Sub(order.LiveStartDate).Seconds() / 86400)

	order.HomestayTotalPrice = int64(resp.Homestay.HomestayPrice * liveDays)
	if in.IsFood {
		order.NeedFood = model.HomestayOrderNeedFoodYes
		order.FoodTotalPrice = int64(resp.Homestay.FoodPrice * in.LivePeopleNum * liveDays)
	}

	_, err = l.svcCtx.HomestayOrderModel.Insert(l.ctx, nil, order)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Order Database Exception order : %+v , err: %v", order, err)
	}

	//设置订单定时关闭
	//todo...

	return &pb.CreateHomestayOrderResp{
		Sn: order.Sn,
	}, nil
}
