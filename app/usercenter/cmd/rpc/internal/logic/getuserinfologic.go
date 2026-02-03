package logic

import (
	"context"

	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/pb"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"

	"github.com/jinzhu/copier"
)

var ErrUserNoExistsError = xerr.NewErrMsg("用户不存在")

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "GetUserInfo find user db err , id:%d , err:%v", in.Id, err)
	}
	if user == nil {
		return nil, errors.Wrapf(ErrUserNoExistsError, "id:%d", in.Id)
	}

	var respUser pb.User
	_ = copier.Copy(&respUser, user)

	return &pb.GetUserInfoResp{
		User: &respUser,
	}, nil
}
