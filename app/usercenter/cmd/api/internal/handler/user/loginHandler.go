// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"net/http"

	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/logic/user"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/pkg/result"
	"github.com/zeromicro/go-zero/rest/httpx"
)

// login
func LoginHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.LoginReq
		if err := httpx.Parse(r, &req); err != nil {
			result.ParamErrorResult(r, w, err)
			return
		}

		l := user.NewLoginLogic(r.Context(), svcCtx)
		resp, err := l.Login(&req)
		result.HttpResult(r, w, resp, err)
	}
}
