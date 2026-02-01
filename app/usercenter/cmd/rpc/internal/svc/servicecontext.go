package svc

import "github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/internal/config"

type ServiceContext struct {
	Config config.Config
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,
	}
}
