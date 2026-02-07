// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package homestayBussiness

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/travel/model"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/usercenter"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type GoodBossLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 推荐房主
func NewGoodBossLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GoodBossLogic {
	return &GoodBossLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GoodBossLogic) GoodBoss(req *types.GoodBossReq) (*types.GoodBossResp, error) {
	//todo... Mock实现,选择策略待补全
	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(
		squirrel.Eq{
			"row_type":   model.HomestayActivityGoodBusiType,
			"row_status": model.HomestayActivityUpStatus,
		},
	)
	hometstayActivityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx,
		whereBuilder,
		0,
		10,
		"data_id desc",
	)
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "get GoodBoss db err. rowType: %s ,err : %v", model.HomestayActivityGoodBusiType, err)
	}

	var resp []types.HomestayBusinessBoss
	if len(hometstayActivityList) > 0 {

		mr.MapReduceVoid(
			func(source chan<- interface{}) {
				for _, homestayActivity := range hometstayActivityList {
					source <- homestayActivity.DataId
				}
			}, func(item interface{}, writer mr.Writer[*usercenter.User], cancel func(error)) {
				id := item.(int64)

				userResp, err := l.svcCtx.UsercenterRpc.GetUserInfo(l.ctx, &usercenter.GetUserInfoReq{
					Id: id,
				})
				if err != nil {
					logx.WithContext(l.ctx).Errorf("GoodBossLogic GoodBoss fail userId : %d ,err:%v", id, err)
					return
				}

				if userResp.User != nil && userResp.User.Id > 0 {
					writer.Write(userResp.User)
				}
			}, func(pipe <-chan *usercenter.User, cancel func(error)) {
				for item := range pipe {
					var typeHomestayBussiness types.HomestayBusinessBoss
					_ = copier.Copy(&typeHomestayBussiness, item)

					resp = append(resp, typeHomestayBussiness)
				}
			},
		)

	}

	return &types.GoodBossResp{
		List: resp,
	}, nil
}
