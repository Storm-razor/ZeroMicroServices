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

// 猜你喜欢民宿列表
func GuessListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GuessListReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := homestay.NewGuessListLogic(r.Context(), svcCtx)
		resp, err := l.GuessList(&req)
		result.HttpResult(r, w, resp, err)
	}
}
