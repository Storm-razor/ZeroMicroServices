// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package thirdPayment

import (
	"net/http"

	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/logic/thirdPayment"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/payment/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 第三方支付:微信支付
func ThirdPaymentwxPayHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ThirdPaymentWxPayReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := thirdPayment.NewThirdPaymentwxPayLogic(r.Context(), svcCtx)
		resp, err := l.ThirdPaymentwxPay(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
