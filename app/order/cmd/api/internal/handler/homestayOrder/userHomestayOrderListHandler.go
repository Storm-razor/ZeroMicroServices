// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package homestayOrder

import (
	"net/http"

	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/logic/homestayOrder"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/pkg/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 用户订单列表
func UserHomestayOrderListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UserHomestayOrderListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestayOrder.NewUserHomestayOrderListLogic(r.Context(), svcCtx)
		resp, err := l.UserHomestayOrderList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
