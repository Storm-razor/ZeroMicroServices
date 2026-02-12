package logic

import (
	"context"

	"github.com/hibiken/asynq"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/internal/svc"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/jobtype"
)

type CronJob struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCronJob(ctx context.Context, svcCtx *svc.ServiceContext) *CronJob {
	return &CronJob{
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

// ---------------------------
// @brief 注册任务处理回调
// ---------------------------
func (l *CronJob) Register() *asynq.ServeMux {

	mux := asynq.NewServeMux()

	//scheduler job
	mux.Handle(jobtype.ScheduleSettleRecord, NewSettleRecordHandler(l.svcCtx))

	//defer job
	mux.Handle(jobtype.DeferCloseHomestayOrder, NewCloseHomestayOrderHandler(l.svcCtx))

	//SuccessNotifyUser job
	mux.Handle(jobtype.MsgPaySuccessNotifyUser, NewPaySuccessNotifyUserHandler(l.svcCtx))
	//queue job , asynq support queue job

	return mux
}
