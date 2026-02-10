package kq

import (
	"context"
	"encoding/json"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/order"
	"github.com/wwwzy/ZeroMicroServices/pkg/kqueue"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentUpdateStatusMq struct {
	ctx context.Context

	svcCtx *svc.ServiceContext
}

func NewPaymentUpdateStatusMq(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentUpdateStatusMq {
	return &PaymentUpdateStatusMq{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ---------------------------
// @brief 消费者
// ---------------------------
func (l *PaymentUpdateStatusMq) Consume(_ context.Context, _, val string) error {

	var message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->Consume Unmarshal err : %v , val : %s", err, val)
		return err
	}

	if err := l.execService(message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->execService  err : %v , val : %s , message:%+v", err, val, message)
		return err
	}

	return nil
}

// ---------------------------
// @brief 调用rpc更新订单状态
// ---------------------------
func (l *PaymentUpdateStatusMq) execService(message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage) error {

	orderTradeState := l.getOrderTradeStateByPaymentTradeState(message.PayStatus)
	if orderTradeState != -99 {
		_, err := l.svcCtx.OrderRpc.UpdateHomestayOrderTradeState(l.ctx, &order.UpdateHomestayOrderTradeStateReq{
			Sn:         message.OrderSn,
			TradeState: orderTradeState,
		})
		if err != nil {
			return errors.Wrapf(xerr.NewErrMsg("update homestay order state fail"), "update homestay order state fail err : %v ,message:%+v", err, message)
		}
		// RPC调用成功，返回 nil
		return nil
	}

	// 支付状态未定义或无法映射，返回错误提示 TODO
	return errors.Errorf("execService Todo... unknown payment status: %d", message.PayStatus)
}

// ---------------------------
// @brief 根据第三方支付状态获取订单状态
// ---------------------------
func (l *PaymentUpdateStatusMq) getOrderTradeStateByPaymentTradeState(thirdPaymentPayStatus int64) int64 {
	// todo...

	// switch thirdPaymentPayStatus {
	// case paymentModel.ThirdPaymentPayTradeStateSuccess:
	// 	return model.HomestayOrderTradeStateWaitUse
	// case paymentModel.ThirdPaymentPayTradeStateRefund:
	// 	return model.HomestayOrderTradeStateRefund
	// default:
	// 	return -99
	// }

	return -99
}
