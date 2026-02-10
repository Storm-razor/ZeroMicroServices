package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/rpc/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/rpc/pb"
	"github.com/wwwzy/ZeroMicroServices/app/payment/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/kqueue"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateTradeStateLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateTradeStateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateTradeStateLogic {
	return &UpdateTradeStateLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新交易状态
func (l *UpdateTradeStateLogic) UpdateTradeState(in *pb.UpdateTradeStateReq) (*pb.UpdateTradeStateResp, error) {
	//检查支付记录
	thirdPayment, err := l.svcCtx.ThirdPaymentModel.FindOneBySn(l.ctx, in.Sn)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "UpdateTradeState FindOneBySn db err , sn : %s , err : %+v", in.Sn, err)
	}

	if thirdPayment == nil {
		return nil, errors.Wrapf(xerr.NewErrMsg("third payment record no exists"), " sn : %s", in.Sn)
	}

	//检查支付状态
	if in.PayStatus == model.ThirdPaymentPayTradeStateSuccess || in.PayStatus == model.ThirdPaymentPayTradeStateFAIL {
		//若为成功或失败,则检查
		if thirdPayment.PayStatus == model.ThirdPaymentPayTradeStateWait {
			//如果存储的支付记录不是待支付状态,则直接返回
			return &pb.UpdateTradeStateResp{}, nil
		}
	} else if in.PayStatus == model.ThirdPaymentPayTradeStateRefund {
		//若退款,则检查
		if thirdPayment.PayStatus != model.ThirdPaymentPayTradeStateSuccess {
			//如果存储的支付记录不是成功状态,则返回错误
			return nil, errors.Wrapf(xerr.NewErrMsg("Only orders with successful payment can be refunded"), "Only orders with successful payment can be refunded in : %+v", in)
		}
	} else {
		return nil, errors.Wrapf(xerr.NewErrMsg("This status is not currently supported"), "Modify payment flow status is not supported  in : %+v", in)
	}

	//更新支付数据库
	thirdPayment.TradeState = in.TradeState
	thirdPayment.TransactionId = in.TransactionId
	thirdPayment.TradeType = in.TradeType
	thirdPayment.TradeStateDesc = in.TradeStateDesc
	thirdPayment.PayStatus = in.PayStatus
	thirdPayment.PayTime = time.Unix(in.PayTime, 0)
	if err := l.svcCtx.ThirdPaymentModel.UpdateWithVersion(l.ctx, nil, thirdPayment); err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), " UpdateTradeState UpdateWithVersion db  err:%v ,thirdPayment : %+v , in : %+v", err, thirdPayment, in)
	}

	//将订单状态变化消息放入消息队列
	if err := l.pubkqPaySuccess(in.Sn, in.PayStatus); err != nil {
		logx.WithContext(l.ctx).Errorf("l.pubKqPaySuccess : %+v", err)
	}

	return &pb.UpdateTradeStateResp{}, nil
}

// ---------------------------
// @brief 向kq队列中放入消息
// ---------------------------
func (l *UpdateTradeStateLogic) pubkqPaySuccess(orderSn string, payStatus int64) error {
	m := kqueue.ThirdPaymentUpdatePayStatusNotifyMessage{
		OrderSn:   orderSn,
		PayStatus: payStatus,
	}

	body, err := json.Marshal(m)
	if err != nil {
		return errors.Wrapf(xerr.NewErrMsg("kq UpdateTradeStateLogic pushKqPaySuccess task marshal error "), "kq UpdateTradeStateLogic pushKqPaySuccess task marshal error  , v : %+v", m)
	}

	return l.svcCtx.KqueuePaymentUpdatePayStatusClient.Push(l.ctx, string(body))
}
