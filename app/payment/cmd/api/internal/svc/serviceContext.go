// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/order"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/rpc/payment"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/usercenter"
	"github.com/zeromicro/go-zero/zrpc"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
)

type ServiceContext struct {
	Config config.Config

	WxPayClient *core.Client

	PaymentRpc    payment.Payment
	OrderRpc      order.Order
	UsercenterRpc usercenter.Usercenter
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config: c,

		PaymentRpc:    payment.NewPayment(zrpc.MustNewClient(c.PaymentRpcConf)),
		OrderRpc:      order.NewOrder(zrpc.MustNewClient(c.OrderRpcConf)),
		UsercenterRpc: usercenter.NewUsercenter(zrpc.MustNewClient(c.UsercenterRpcConf)),
	}
}
