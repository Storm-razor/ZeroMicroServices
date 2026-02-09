package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/pb"
	"github.com/wwwzy/ZeroMicroServices/app/order/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateHomestayOrderTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateHomestayOrderTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateHomestayOrderTradeStateLogic {
	return &UpdateHomestayOrderTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新民宿订单状态服务
func (l *UpdateHomestayOrderTradeStateLogic) UpdateHomestayOrderTradeState(in *pb.UpdateHomestayOrderTradeStateReq) (*pb.UpdateHomestayOrderTradeStateResp, error) {
	// 检查当前订单
	homestayOrder, err := l.svcCtx.HomestayOrderModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "UpdateHomestayOrderTradeState FindOneBySn db err : %v , in:%+v", err, in)
	}
	if homestayOrder == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("order no exists"), "order no exists  in : %+v", in)
	}

	if homestayOrder.TradeState == in.TradeState {
		return &pb.UpdateHomestayOrderTradeStateResp{}, nil
	}

	// 验证状态是否可更改
	if err := l.verifyOrderTradeState(in.TradeState, homestayOrder.TradeState); err != nil {
		return nil, errors.WithMessagef(err, " , in : %+v", in)
	}

	// 更新订单状态
	homestayOrder.TradeState = in.TradeState
	if err := l.svcCtx.HomestayOrderModel.UpdateWithVersion(l.ctx, nil, homestayOrder); err != nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("Failed to update homestay order status"), "Failed to update homestay order status db UpdateWithVersion err:%v , in : %v", err, in)
	}

	// 在队列中设置用户提示
	// todo...

	return &pb.UpdateHomestayOrderTradeStateResp{
		Id:              homestayOrder.Id,
		UserId:          homestayOrder.UserId,
		Sn:              homestayOrder.Sn,
		TradeCode:       homestayOrder.TradeCode,
		Title:           homestayOrder.Title,
		LiveStartDate:   homestayOrder.LiveStartDate.Unix(),
		LiveEndDate:     homestayOrder.LiveEndDate.Unix(),
		OrderTotalPrice: homestayOrder.OrderTotalPrice,
	}, nil

}

func (l *UpdateHomestayOrderTradeStateLogic) verifyOrderTradeState(newTradeState, oldTradeState int64) error {
	if newTradeState == model.HomestayOrderTradeStateWaitPay {
		return errors.Wrapf(xerr.NewErrMsg("Changing this status is not supported"),
			"Changing this status is not supported newTradeState: %d, oldTradeState: %d",
			newTradeState,
			oldTradeState)
	}

	if newTradeState == model.HomestayOrderTradeStateCancel {
		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return errors.Wrapf(xerr.NewErrMsg("只有待支付的订单才能被取消"),
				"Only orders pending payment can be cancelled newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}

	} else if newTradeState == model.HomestayOrderTradeStateWaitUse {
		if oldTradeState != model.HomestayOrderTradeStateWaitPay {
			return errors.Wrapf(xerr.NewErrMsg("Only orders pending payment can change this status"),
				"Only orders pending payment can change this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateUsed {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateRefund {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	} else if newTradeState == model.HomestayOrderTradeStateExpire {
		if oldTradeState != model.HomestayOrderTradeStateWaitUse {
			return errors.Wrapf(xerr.NewErrMsg("Only unused orders can be changed to this status"),
				"Only unused orders can be changed to this status newTradeState: %d, oldTradeState: %d",
				newTradeState,
				oldTradeState)
		}
	}
	return nil
}
