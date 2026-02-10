// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package homestayOrder

import (
	"net/http"

	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/logic/homestayOrder"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户订单明细
func UserHomestayOrderDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserHomestayOrderDetailReq
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := homestayOrder.NewUserHomestayOrderDetailLogic(r.Context(), svcCtx)
		resp, err := l.UserHomestayOrderDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
