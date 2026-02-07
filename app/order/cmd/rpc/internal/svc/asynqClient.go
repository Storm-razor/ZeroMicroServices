package svc

import (
	"github.com/hibiken/asynq"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/internal/config"
)

// 创建AsynqClient
func newAsynqClient(c config.Config) *asynq.Client {
	return asynq.NewClient(asynq.RedisClientOpt{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
	})
}
