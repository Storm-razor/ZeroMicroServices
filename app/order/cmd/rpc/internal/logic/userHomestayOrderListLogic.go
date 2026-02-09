package logic

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/order/cmd/rpc/pb"
	"github.com/wwwzy/ZeroMicroServices/app/order/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UserHomestayOrderListLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUserHomestayOrderListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UserHomestayOrderListLogic {
	return &UserHomestayOrderListLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户民宿订单服务
func (l *UserHomestayOrderListLogic) UserHomestayOrderList(in *pb.UserHomestayOrderListReq) (*pb.UserHomestayOrderListResp, error) {
	// 构建查询条件：固定筛选当前用户
	whereBuilder := l.svcCtx.HomestayOrderModel.SelectBuilder().Where(squirrel.Eq{
		"user_id": in.UserId,
	})
	// 动态筛选：如果前端传入了特定的订单状态，进一步筛选
	if in.TraderState >= model.HomestayOrderTradeStateCancel && in.TraderState <= model.HomestayOrderTradeStateExpire {
		whereBuilder = whereBuilder.Where(squirrel.Eq{"trade_state": in.TraderState})
	}

	list, err := l.svcCtx.HomestayOrderModel.FindPageListByIdDESC(l.ctx,
		whereBuilder,
		in.LastId,
		in.PageSize,
	)
	if err != nil && err != model.ErrNotFound {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "Failed to get user's homestay order err : %v , in :%+v", err, in)
	}

	var resp []*pb.HomestayOrder
	if len(list) > 0 {
		for _, item := range list {
			var pbHomestayOrder pb.HomestayOrder
			_ = copier.Copy(&pbHomestayOrder, item)

			pbHomestayOrder.CreateTime = item.CreateTime.Unix()
			pbHomestayOrder.LiveStartDate = item.LiveStartDate.Unix()
			pbHomestayOrder.LiveEndDate = item.LiveEndDate.Unix()

			resp = append(resp, &pbHomestayOrder)
		}
	}

	return &pb.UserHomestayOrderListResp{
		List: resp,
	}, nil
}
