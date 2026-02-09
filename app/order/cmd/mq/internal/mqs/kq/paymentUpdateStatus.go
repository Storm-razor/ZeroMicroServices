package kq

import (
	"context"
	"encoding/json"

	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/pkg/kqueue"
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
func (l *PaymentUpdateStatusMq) Consume(_, val string) error {

	var message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage
	if err := json.Unmarshal([]byte(val), &message); err != nil {
		logx.WithContext(l.ctx).Error("PaymentUpdateStatusMq->Consume Unmarshal err : %v , val : %s", err, val)
		return err
	}

	// todo...

	return nil
}

// ---------------------------
// @brief 调用rpc更新订单状态
// ---------------------------
func (l *PaymentUpdateStatusMq) execService(message kqueue.ThirdPaymentUpdatePayStatusNotifyMessage) error {
	// todo...
	return nil
}
