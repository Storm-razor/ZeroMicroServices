// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/order"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/rpc/payment"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/rpc/travel"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config config.Config

	OrderRpc   order.Order
	PaymentRpc payment.Payment
	TravelRpc  travel.Travel
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		OrderRpc:   order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		PaymentRpc: payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		TravelRpc:  travel.NewTravel(zrpc.MustNewClient(c.TravelRpcConf)),
	}
}
