package listen

import (
	"context"

	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/mq/internal/svc"

	"github.com/zeromicro/go-zero/core/service"
)

// ---------------------------
// @brief 创建消费者服务
// ---------------------------
func Mqs(c config.Config) []service.Service {
	svcContext := svc.NewServieContext(c)
	ctx := context.Background()

	var services []service.Service

	services = append(services, KqMqs(c, ctx, svcContext)...)

	return services
}
