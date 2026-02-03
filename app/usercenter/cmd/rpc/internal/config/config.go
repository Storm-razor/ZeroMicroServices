package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	zrpc.RpcServerConf
	JwtAuth struct {
		AccessSecret string //加密/解密密钥(服务端保密)
		AccessExpire int64  //token的有效时长
	}
	DB struct {
		DataSource string
	}
	Cache cache.CacheConf
}
