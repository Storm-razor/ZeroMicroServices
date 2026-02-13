// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package homestay

import (
	"net/http"

	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/logic/homestay"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/pkg/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// 民宿列表
func HomestayListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.HomestayListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewHomestayListLogic(r.Context(), svcCtx)
		resp, err := l.HomestayList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
