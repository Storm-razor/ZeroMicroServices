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
			// 增加下面这两行，确保错误时直接结束，并把 "FAIL" 传给微信
			if resp != nil {
				fmt.Fprint(w, resp.ReturnCode)
			}
			return
		}

		w.WriteHeader(http.StatusOK)
		if resp != nil {
			fmt.Fprint(w, resp.ReturnCode)
		}
	}
}
