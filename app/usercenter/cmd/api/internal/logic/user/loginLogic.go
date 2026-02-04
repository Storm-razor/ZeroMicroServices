// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"

	"github.com/jinzhu/copier"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/usercenter"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (*types.LoginResp, error) {
	loginResp, err := l.svcCtx.UserCenterRpc.Login(l.ctx, &usercenter.LoginReq{
		AuthType: model.UserAuthTypeSystem,
		AuthKey:  req.Mobile,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	var resp types.LoginResp
	_ = copier.Copy(&resp, loginResp)

	return &resp, nil
}
