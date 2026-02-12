package logic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/golang-module/carbon/v2"
	"github.com/hibiken/asynq"
	"github.com/pkg/errors"
	"github.com/silenceper/wechat/v2/miniprogram/subscribe"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/jobtype"
	"github.com/wwwzy/ZeroMicroServices/app/order/model"
	"github.com/wwwzy/ZeroMicroServices/app/usercenter/cmd/rpc/usercenter"
	usercenterModel "github.com/wwwzy/ZeroMicroServices/app/usercenter/model"
	"github.com/wwwzy/ZeroMicroServices/pkg/globalkey"
	"github.com/wwwzy/ZeroMicroServices/pkg/tool"
	"github.com/wwwzy/ZeroMicroServices/pkg/wechat"
	"github.com/wwwzy/ZeroMicroServices/pkg/xerr"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/service"
)

var ErrPaySuccessNotifyFail = xerr.NewErrMsg("pay success notify user fail")

type PaySuccessNotifyUserHandler struct {
	svcCtx *svc.ServiceContext
}

func NewPaySuccessNotifyUserHandler(svcCtx *svc.ServiceContext) *PaySuccessNotifyUserHandler {
	return &PaySuccessNotifyUserHandler{
		svcCtx: svcCtx,
	}
}

// ---------------------------
// @brief 提醒用户支付状态
// ---------------------------
func (l *PaySuccessNotifyUserHandler) ProcessTask(ctx context.Context, t *asynq.Task) error {
	// 解析任务参数
	var p jobtype.PaySuccessNotifyUserPayload
	if err := json.Unmarshal(t.Payload(), &p); err != nil {
		return errors.Wrapf(ErrCloseOrderFal, "closeHomestayOrderStateMqHandler payload err:%v, payLoad:%+v", err, t.Payload())
	}

	// 获取用户微信登陆ID
	usercenterResp, err := l.svcCtx.UsercenterRpc.GetUserAuthByUserId(ctx, &usercenter.GetUserAuthByUserIdReq{
		UserId:   p.Order.UserId,
		AuthType: usercenterModel.UserAuthTypeSmallWX,
	})
	if err != nil {
		return errors.Wrapf(ErrPaySuccessNotifyFail, "pay success notify user fail, rpc get user err:%v , orderSn:%s , userId:%d", err, p.Order.Sn, p.Order.UserId)
	}
	if usercenterResp.UserAuth == nil || len(usercenterResp.UserAuth.AuthKey) == 0 {
		return errors.Wrapf(ErrPaySuccessNotifyFail, "pay success notify user , user no exists err:%v , orderSn:%s , userId:%d", err, p.Order.Sn, p.Order.UserId)
	}
	openId := usercenterResp.UserAuth.AuthKey

	//发送信息
	msgs := l.getData(ctx, p.Order, openId)
	for _, msg := range msgs {
		l.SendWxMini(ctx, msg)
	}

	return nil
}

// ---------------------------
// @brief 生成用户消息
// ---------------------------
func (l *PaySuccessNotifyUserHandler) getData(_ context.Context, order *model.HomestayOrder, openId string) []*subscribe.Message {
	return []*subscribe.Message{
		{
			ToUser:     openId,
			TemplateID: wechat.OrderPaySuccessTemplateID,
			Data: map[string]*subscribe.DataItem{
				"character_string6": {Value: order.Sn},
				"thing1":            {Value: order.Title},
				"amount2":           {Value: fmt.Sprintf("%.2f", tool.Fen2Yuan(order.OrderTotalPrice))},
				"time4":             {Value: carbon.CreateFromTimestamp(order.LiveStartDate.Unix()).Format(globalkey.DateTimeFormatTplStandardDate)},
				"time5":             {Value: carbon.CreateFromTimestamp(order.LiveEndDate.Unix()).Format(globalkey.DateTimeFormatTplStandardDate)},
			},
		},
		{
			ToUser:     openId,
			TemplateID: wechat.OrderPaySuccessLiveKnowTemplateID,
			Data: map[string]*subscribe.DataItem{
				"date2":             {Value: carbon.CreateFromTimestamp(order.LiveStartDate.Unix()).Format(globalkey.DateTimeFormatTplStandardDate)},
				"date3":             {Value: carbon.CreateFromTimestamp(order.LiveEndDate.Unix()).Format(globalkey.DateTimeFormatTplStandardDate)},
				"character_string4": {Value: order.TradeCode},
				"thing1":            {Value: "请不要将验证码告知商家以外人员，以防上当"},
			},
		},
	}

}

// ---------------------------
// @brief 通过微信发送消息
// ---------------------------
func (l *PaySuccessNotifyUserHandler) SendWxMini(ctx context.Context, msg *subscribe.Message) {
	if l.svcCtx.Config.Mode != service.PreMode {
		msg.MiniprogramState = "developer"
	} else {
		msg.MiniprogramState = "formal"
	}

	var maxRetryNum int64 = 5
	var retryNum int64

	for {
		time.Sleep(time.Second)
		err := l.svcCtx.MiniProgram.GetSubscribe().Send(msg)
		if err != nil {
			if retryNum > maxRetryNum {
				logx.WithContext(ctx).Errorf("Payment successful send wechat mini subscription message failed retryNum ： %d , err:%v, msg ： %+v ", retryNum, err, msg)
				return
			}
			retryNum++
			continue
		}
		return
	}
}
