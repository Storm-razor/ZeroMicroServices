package svc

import (
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/model"

	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config

	RedisClient   *redis.Redis
	UserModel     model.UserModel
	UserAuthModel model.UserAuthModel
}

func NewServiceContext(c config.Config) *ServiceContext {
	sqlConn := sqlx.NewMysql(c.DB.DataSource)

	return &ServiceContext{
		Config:        c,
		UserModel:     model.NewUserModel(sqlConn),
		UserAuthModel: model.NewUserAuthModel(sqlConn),
		RedisClient: redis.New(c.Redis.Host, func(r *redis.Redis) {
			r.Type = c.Redis.Type
			r.Pass = c.Redis.Pass
		}),
	}
}
