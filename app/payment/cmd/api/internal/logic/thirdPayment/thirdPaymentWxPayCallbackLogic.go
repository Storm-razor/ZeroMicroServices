// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package thirdPayment

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/rpc/payment"
	"github.com/wwwzy/ZeroMicroServices/app/payment/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/wechatpay-apiv3/wechatpay-go/core/auth/verifiers"
	"github.com/wechatpay-apiv3/wechatpay-go/core/downloader"
	"github.com/wechatpay-apiv3/wechatpay-go/core/notify"
	"github.com/wechatpay-apiv3/wechatpay-go/services/payments"

	"github.com/zeromicro/go-zero/core/logx"
)

var ErrWxPayCallbackError = xerr.NewErrMsg("wechat pay callback fail")

type ThirdPaymentWxPayCallbackLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 第三方支付:微信支付回调
func NewThirdPaymentWxPayCallbackLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ThirdPaymentWxPayCallbackLogic {
	return &ThirdPaymentWxPayCallbackLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ---------------------------
// @brief 微信支付回调必须读取原始http包体,所以需要自定义处理逻辑
// ---------------------------
func (l *ThirdPaymentWxPayCallbackLogic) ThirdPaymentWxPayCallback(_ http.ResponseWriter, req *http.Request) (*types.ThirdPaymentWxPayCallbackResp, error) {

	_, err := svc.NewWxPayClientV3(l.svcCtx.Config)
	if err != nil {
		return nil, err
	}

	certVisitor := downloader.MgrInstance().GetCertificateVisitor(l.svcCtx.Config.WxPayConf.MchId)
	handler := notify.NewNotifyHandler(l.svcCtx.Config.WxPayConf.APIv3Key, verifiers.NewSHA256WithRSAVerifier(certVisitor))

	transaction := new(payments.Transaction)
	_, err = handler.ParseNotifyRequest(context.Background(), req, transaction)
	if err != nil {
		return nil, errors.Wrapf(ErrWxPayCallbackError, "Failed to parse data ,err:%v", err)
	}

	returnCode := "SUCCESS"
	err = l.verifyAndUpdateState(transaction)
	if err != nil {
		returnCode = "FAIL"
	}

	return &types.ThirdPaymentWxPayCallbackResp{
		ReturnCode: returnCode,
	}, err

}

// ---------------------------
// @brief 验证并更新支付状态
// ---------------------------
func (l *ThirdPaymentWxPayCallbackLogic) verifyAndUpdateState(notifyTrasaction *payments.Transaction) error {
	// 根据支付号查询支付记录
	payementResp, err := l.svcCtx.PaymentRpc.GetPaymentBySn(l.ctx, &payment.GetPaymentBySnReq{
		Sn: *notifyTrasaction.OutTradeNo,
	})
	if err != nil || payementResp.PaymentDetail.Id == 0 {
		return errors.Wrapf(ErrWxPayCallbackError, "Failed to get payment flow record err:%v ,notifyTrasaction:%+v ", err, notifyTrasaction)
	}

	// 对比金额
	notifyPayTotal := *notifyTrasaction.Amount.PayerTotal
	if payementResp.PaymentDetail.PayTotal != notifyPayTotal {
		return errors.Wrapf(ErrWxPayCallbackError, "Order amount exception  notifyPayTotal:%v , notifyTrasaction:%v ", notifyPayTotal, notifyTrasaction)
	}

	// 检查支付状态
	payStatus := l.getPayStatusByWXPayTradeState(*notifyTrasaction.TradeState)
	if payStatus == model.ThirdPaymentPayTradeStateSuccess {
		if payementResp.PaymentDetail.PayStatus != model.ThirdPaymentPayTradeStateWait {
			return nil
		}
		//更新状态
		if _, err := l.svcCtx.PaymentRpc.UpdateTradeState(l.ctx, &payment.UpdateTradeStateReq{
			Sn:             *notifyTrasaction.OutTradeNo,
			TradeState:     *notifyTrasaction.TradeState,
			TransactionId:  *notifyTrasaction.TransactionId,
			TradeType:      *notifyTrasaction.TradeType,
			TradeStateDesc: *notifyTrasaction.TradeStateDesc,
			PayStatus:      l.getPayStatusByWXPayTradeState(*notifyTrasaction.TradeState),
		}); err != nil {
			return errors.Wrapf(ErrWxPayCallbackError, "更新流水状态失败  err:%v , notifyTrasaction:%v ", err, notifyTrasaction)
		}

	} else if payStatus == model.ThirdPaymentPayTradeStateWait {
		// 重新发起支付,待实现
		// todo...
	}

	return nil
}

const (
	SUCCESS    = "SUCCESS"    //支付成功
	REFUND     = "REFUND"     //转入退款
	NOTPAY     = "NOTPAY"     //未支付
	CLOSED     = "CLOSED"     //已关闭
	REVOKED    = "REVOKED"    //已撤销（付款码支付）
	USERPAYING = "USERPAYING" //用户支付中（付款码支付）
	PAYERROR   = "PAYERROR"   //支付失败(其他原因，如银行返回失败)
)

// ---------------------------
// @brief 将微信支付服务状态转为后端支付状态
// ---------------------------
func (l *ThirdPaymentWxPayCallbackLogic) getPayStatusByWXPayTradeState(wxPayTradeState string) int64 {
	switch wxPayTradeState {
	case SUCCESS: //支付成功
		return model.ThirdPaymentPayTradeStateSuccess
	case USERPAYING: //支付中
		return model.ThirdPaymentPayTradeStateWait
	case REFUND: //已退款
		return model.ThirdPaymentPayTradeStateRefund
	default:
		return model.ThirdPaymentPayTradeStateFAIL
	}
}
