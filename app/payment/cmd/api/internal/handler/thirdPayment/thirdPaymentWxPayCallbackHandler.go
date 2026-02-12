// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package thirdPayment

import (
	"fmt"
	"net/http"

	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/logic/thirdPayment"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

// 第三方支付:微信支付回调(自定义)
func ThirdPaymentWxPayCallbackHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		l := thirdPayment.NewThirdPaymentWxPayCallbackLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentWxPayCallback(w, r)

		if err != nil {
			logx.WithContext(r.Context()).Errorf("【API-ERR】 ThirdPaymentWxPayCallbackHandler : %+v ", err)
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}

		logx.Infof("ReturnCode : %s ", resp.ReturnCode)
		fmt.Fprint(w.(http.ResponseWriter), resp.ReturnCode)
		// var req types.ThirdPaymentWxPayCallbackReq
		// if err := httpx.Parse(r, &req); err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// 	return
		// }

		// l := thirdPayment.NewThirdPaymentWxPayCallbackLogic(r.Context(), svcCtx)
		// resp, err := l.ThirdPaymentWxPayCallback(&req)
		// if err != nil {
		// 	httpx.ErrorCtx(r.Context(), w, err)
		// } else {
		// 	httpx.OkJsonCtx(r.Context(), w, resp)
		// }
	}
}
