// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package thirdPayment

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/order"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/rpc/payment"
	"github.com/wwwzy/ZeroMicroServices/app/payment/model"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/usercenter"
	usercentermodel "github.com/wwwzy/ZeroMicroServices/app/usercenter/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/ctxdata"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/wechatpay-apiv3/wechatpay-go/core"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments/jsapi"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxPayError = xerr.NewErrMsg("wechat pay fail")

type ThirdPaymentwxPayLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 第三方支付:微信支付
func NewThirdPaymentwxPayLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentwxPayLogic {
	return &ThirdPaymentwxPayLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ThirdPaymentwxPayLogic) ThirdPaymentwxPay(req *types.ThirdPaymentWxPayReq) (*types.ThirdPaymentWxPayResp, error) {
	var totalPrice int64   // 支付金额，单位为分
	var description string // 支付订单描述

	//获取订单总价格和信息
	switch req.ServiceType {
	case model.ThirdPaymentServiceTypeHomestayOrder:
		homestayTotalPrice, homestayDescription, err := l.getPayHomestayPriceDescription(req.OrderSn)
		if err != nil {
			return nil, errors.Wrapf(ErrWxPayError, "getPayHomestayPriceDescription err : %v req: %+v", err, req)
		}
		totalPrice = homestayTotalPrice
		description = homestayDescription
	default:
		return nil, errors.Wrapf(xerr.NewErrMsg("Payment for this business type is not supported"), "Payment for this business type is not supported req: %+v", req)
	}

	//调用微信支付接口
	wechatPrepayRsp, err := l.createWxPrePayOrder(req.ServiceType, req.OrderSn, totalPrice, description)
	if err != nil {
		return nil, err
	}

	return &types.ThirdPaymentWxPayResp{
		Appid:     l.svcCtx.Config.WxMiniConf.AppId,
		NonceStr:  *wechatPrepayRsp.NonceStr,
		PaySign:   *wechatPrepayRsp.PaySign,
		Package:   *wechatPrepayRsp.Package,
		Timestamp: *wechatPrepayRsp.TimeStamp,
		SignType:  *wechatPrepayRsp.SignType,
	}, nil
}

// ---------------------------
// @brief 根据订单号获取总支付金额和订单描述
// ---------------------------
func (l *ThirdPaymentwxPayLogic) getPayHomestayPriceDescription(orderSn string) (int64, string, error) {

	description := "homestay pay"

	resp, err := l.svcCtx.OrderRpc.HomestayOrderDetail(l.ctx, &order.HomestayOrderDetailReq{
		Sn: orderSn,
	})
	if err != nil {
		return 0, description, errors.Wrapf(ErrWxPayError,
			"OrderRpc.HomestayOrderDetail err: %v, orderSn: %s", err, orderSn)
	}
	if resp.HomestayOrder == nil || resp.HomestayOrder.Id == 0 {
		return 0, description, errors.Wrapf(xerr.NewErrMsg("order no exists"), "WeChat payment order does not exist orderSn : %s", orderSn)
	}

	return resp.HomestayOrder.OrderTotalPrice, description, nil
}

// ---------------------------
// @brief 调用微信支付接口
// ---------------------------
func (l *ThirdPaymentwxPayLogic) createWxPrePayOrder(serviceType, orderSn string, totalPrice int64, description string) (*jsapi.PrepayWithRequestPaymentResponse, error) {

	//通过用户id获取微信登陆的id
	userId := ctxdata.GetUidFromCtx(l.ctx)
	userResp, err := l.svcCtx.UsercenterRpc.GetUserAuthByUserId(l.ctx, &usercenter.GetUserAuthByUserIdReq{
		UserId:   userId,
		AuthType: usercentermodel.UserAuthTypeSmallWX,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrWxPayError, "Get user wechat openid err : %v , userId: %d , orderSn:%s", err, userId, orderSn)
	}
	if userResp.UserAuth == nil || userResp.UserAuth.Id == 0 {
		return nil, errors.Wrapf(xerr.NewErrMsg("Get user wechat openid fail，Please pay before authorization by weChat"), "Get user WeChat openid does not exist  userId: %d , orderSn:%s", userId, orderSn)
	}
	openId := userResp.UserAuth.AuthKey

	// 通过paymentRPC创建支付记录
	createPaymentResp, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentReq{
		UserId:      userId,
		PayModel:    model.ThirdPaymentPayModelWechatPay,
		PayTotal:    totalPrice,
		OrderSn:     orderSn,
		ServiceType: serviceType,
	})
	if err != nil || createPaymentResp.Sn == "" {
		return nil, errors.Wrapf(ErrWxPayError,
			"create local third payment record fail : err: %v , userId: %d,totalPrice: %d , orderSn: %s",
			err, userId, totalPrice, orderSn)
	}

	// 创建微信支付服务客户端
	wxPayClient, err := svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}
	jsApiSvc := jsapi.JsapiApiService{Client: wxPayClient}

	// 调用微信支付接口
	resp, _, err := jsApiSvc.PrepayWithRequestPayment(l.ctx,
		jsapi.PrepayRequest{
			Appid:       core.String(l.svcCtx.Config.WxMiniConf.AppId),
			Mchid:       core.String(l.svcCtx.Config.WxPayConf.MchId),
			Description: core.String(description),
			OutTradeNo:  core.String(createPaymentResp.Sn),
			Attach:      core.String(description),
			NotifyUrl:   core.String(l.svcCtx.Config.WxPayConf.NotifyUrl),
			Amount: &jsapi.Amount{
				Total: core.Int64(totalPrice),
			},
			Payer: &jsapi.Payer{
				Openid: core.String(openId),
			},
		},
	)
	if err != nil {
		return nil, errors.Wrapf(ErrWxPayError, "Failed to initiate WeChat payment pre-order err : %v , userId: %d , orderSn:%s", err, userId, orderSn)
	}

	return resp, nil
}
