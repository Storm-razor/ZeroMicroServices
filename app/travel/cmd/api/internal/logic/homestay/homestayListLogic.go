// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package homestay

import (
	"context"

	"github.com/Masterminds/squirrel"
	"github.com/jinzhu/copier"
	"github.com/pkg/errors"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/travel/cmd/api/internal/types"
	"github.com/wwwzy/ZeroMicroServices/app/travel/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/tool"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/mr"
)

type HomestayListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 民宿列表
func NewHomestayListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HomestayListLogic {
	return &HomestayListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HomestayListLogic) HomestayList(req *types.HomestayListReq) (*types.HomestayListResp, error) {
	whereBuilder := l.svcCtx.HomestayActivityModel.SelectBuilder().Where(
		squirrel.Eq{
			"row_type":   model.HomestayActivityPreferredType,
			"row_status": model.HomestayActivityUpStatus,
		},
	)

	//根据民宿活动表查到符合条件的民宿id
	homestayActicityList, err := l.svcCtx.HomestayActivityModel.FindPageListByPage(l.ctx, whereBuilder, req.Page, req.PageSize, "data_id desc")
	if err != nil {
		return nil, errors.Wrapf(xerr.NewErrCode(xerr.DB_ERROR), "get activity homestay id set fail rowType: %s ,err : %v", model.HomestayActivityPreferredType, err)
	}

	var resp []types.Homestay
	if len(homestayActicityList) > 0 {
		mr.MapReduceVoid(func(source chan<- interface{}) {
			for _, homestayActivity := range homestayActicityList {
				source <- homestayActivity.DataId
			}
		}, func(item interface{}, writer mr.Writer[*model.Homestay], cancel func(error)) {
			id := item.(int64)

			homestay, err := l.svcCtx.HomestayModel.FindOne(l.ctx, id)
			if err != nil && err != model.ErrNotFound {
				logx.WithContext(l.ctx).Errorf("ActivityHomestayListLogic ActivityHomestayList 获取活动数据失败 id : %d ,err : %v", id, err)
				return
			}
			writer.Write(homestay)
		}, func(pipe <-chan *model.Homestay, cancel func(error)) {
			for homestay := range pipe {
				var typeHomestay types.Homestay

				_ = copier.Copy(&typeHomestay, homestay)

				typeHomestay.FoodPrice = tool.Fen2Yuan(homestay.FoodPrice)
				typeHomestay.HomestayPrice = tool.Fen2Yuan(homestay.HomestayPrice)
				typeHomestay.MarketHomestayPrice = tool.Fen2Yuan(homestay.MarketHomestayPrice)

				resp = append(resp, typeHomestay)
			}
		})
	}

	return &types.HomestayListResp{
		List: resp,
	}, nil

}
