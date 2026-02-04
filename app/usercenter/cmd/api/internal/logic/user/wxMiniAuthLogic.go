// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package user

import (
	"context"
	"fmt"

	"github.com/pkg/errors"

	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/usercenter"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	miniConfig "github.com/silenceper/wechat/v2/miniprogram/config"
)

var ErrWxMiniAuthFailError = xerr.NewErrMsg("wechat mini auth fail")

type WxMiniAuthLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// wechat mini auth
func NewWxMiniAuthLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WxMiniAuthLogic {
	return &WxMiniAuthLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ---------------------------
// @brief 微信小程序登陆验证
// ---------------------------
func (l *WxMiniAuthLogic) WxMiniAuth(req *types.WXMiniAuthReq) (*types.WXMiniAuthResp, error) {

	//解密微信小程序登陆请求
	miniprogram := wechat.NewWechat().GetMiniProgram(&miniConfig.Config{
		AppID:     l.svcCtx.Config.WxMiniConf.AppId,
		AppSecret: l.svcCtx.Config.WxMiniConf.Secret,
		Cache:     cache.NewMemory(), //使用内存缓存,后期可修改为redis缓存
	})
	authResult, err := miniprogram.GetAuth().Code2Session(req.Code)
	if err != nil || authResult.ErrCode != 0 || authResult.OpenID == "" {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "发起授权请求失败 err : %v , code : %s  , authResult : %+v", err, req.Code, authResult)
	}

	//解析用户数据
	userData, err := miniprogram.GetEncryptor().Decrypt(authResult.SessionKey, req.EncryptedData, req.IV)
	if err != nil {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "解析数据失败 req : %+v , err: %v , authResult:%+v ", req, err, authResult)
	}

	//绑定到用户或者登陆
	var userId int64
	rpcRsp, err := l.svcCtx.UserCenterRpc.GetUserAuthByAuthKey(l.ctx, &usercenter.GetUserAuthByAuthKeyReq{
		AuthType: model.UserAuthTypeSmallWX,
		AuthKey:  authResult.OpenID,
	})
	if err != nil {
		return nil, errors.Wrapf(ErrWxMiniAuthFailError, "rpc call userAuthByAuthKey err : %v , authResult : %+v", err, authResult)
	}
	if rpcRsp.UserAuth == nil || rpcRsp.UserAuth.UserId == 0 {
		//若不存在对应用户则注册
		mobile := userData.PhoneNumber
		nickName := fmt.Sprintf("ZeroMS%s", mobile[7:])
		reigisterRsp, err := l.svcCtx.UserCenterRpc.Register(l.ctx, &usercenter.RegisterReq{
			AuthKey:  authResult.OpenID,
			AuthType: model.UserAuthTypeSmallWX,
			Mobile:   mobile,
			Nickname: nickName,
		})
		if err != nil {
			return nil, errors.Wrapf(ErrWxMiniAuthFailError, "UsercenterRpc.Register err :%v, authResult : %+v", err, authResult)
		}

		return &types.WXMiniAuthResp{
			AccessToken:  reigisterRsp.AccessToken,
			AccessExpire: reigisterRsp.AccessExpire,
			RefreshAfter: reigisterRsp.RefreshAfter,
		}, nil
	} else {
		//若存在对应用户则生成jwt返回
		userId = rpcRsp.UserAuth.UserId
		tokenResp, err := l.svcCtx.UserCenterRpc.GenerateToken(l.ctx, &usercenter.GenerateTokenReq{
			UserId: userId,
		})
		if err != nil {
			return nil, errors.Wrapf(ErrWxMiniAuthFailError, "usercenterRpc.GenerateToken err :%v, userId : %d", err, userId)
		}

		return &types.WXMiniAuthResp{
			AccessToken:  tokenResp.AccessToken,
			AccessExpire: tokenResp.AccessExpire,
			RefreshAfter: tokenResp.RefreshAfter,
		}, nil
	}
}
