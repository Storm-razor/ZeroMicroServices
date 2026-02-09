package config

import (
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	service.ServiceConf

	Redis redis.RedisConf

	PaymentUpdateStatusConf kq.KqConf

	OrderRpcConf      zrpc.RpcClientConf
	UsercenterRpcConf zrpc.RpcClientConf
}
