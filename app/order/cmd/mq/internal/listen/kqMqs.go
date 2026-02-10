package listen

import (
	"context"

	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/config"
	mkq "github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/mqs/kq"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/svc"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
)

// ---------------------------
// @brief 使用gozero的kq库创建消费者队列
// ---------------------------
func KqMqs(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {
	return []service.Service{
		kq.MustNewQueue(c.PaymentUpdateStatusConf, mkq.NewPaymentUpdateStatusMq(ctx, svcContext)),
	}
}
