package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf

	DB struct {
		DataSource string
	}

	Cache cache.CacheConf

	KqPaymentUpdatePayStatusConf KqConfig
}

type KqConfig struct {
	Brokers []string
	Topic   string
}

// 微信小程序配置
type KqServerConfig struct {
	Address string `json:"AppId"`  //微信appId
	Secret  string `json:"Secret"` //微信secret
}
