package svc

import (
	"github.com/hibiken/asynq"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/order/model"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/rpc/travel"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	AsynqClient        *asynq.Client
	HomestayOrderModel model.HomestayOrderModel

	TravelRpc travel.Travel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:             c,
		AsynqClient:        newAsynqClient(c),
		HomestayOrderModel: model.NewHomestayOrderModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		TravelRpc:          travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
