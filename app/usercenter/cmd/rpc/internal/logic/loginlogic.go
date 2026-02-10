package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/pb"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/tool"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

var ErrGenerateTokenError = xerr.NewErrMsg("生成token失败")
var ErrUsernamePwdError = xerr.NewErrMsg("账号或密码不正确")

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	var userId int64
	var err error

	switch in.AuthType {
	case model.UserAuthTypeSystem:
		userId,err = l.loginByMobile(in.AuthKey,in.Password)
	default:
		return nil,xerr.NewErrCode(xerr.SERVER_COMMON_ERROR)
	}
	if err!=nil {
		return nil,err
	}

	// 产生token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx,l.svcCtx)
	tokenResp,err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{
		UserId: userId,
	})
	if err!= nil {
		return nil,errors.Wrapf(ErrGenerateTokenError, "GenerateToken userId : %d", userId)
	}

	return &pb.LoginResp{
		AccessToken:  tokenResp.AccessToken,
		AccessExpire: tokenResp.AccessExpire,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}

func (l *LoginLogic) loginByMobile(mobile, password string) (int64, error) {
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, mobile)
	if err != nil && err != model.ErrNotFound {
		return 0, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "根据手机号查询失败")
	}
	if user == nil {
		return 0, errors.Wrapf(ErrUserNoExistsError, "mobile:%s", mobile)
	}
	if !(tool.Md5ByString(password) == user.Password) {
		return 0, errors.Wrapf(ErrUsernamePwdError, "密码匹配出错")
	}

	return user.Id, nil
}

//---------------------------
//@brief 微信小程序登陆验证 todo...
//---------------------------
func (l *LoginLogic) loginBySmallWx() error {
	return nil
}
