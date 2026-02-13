package logic

import (
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/jobtype"
	"github.com/zeromicro/go-zero/core/logx"
)

func (l *MqueueScheduler) settleRecordScheduler() {
	task := asynq.NewTask(jobtype.ScheduleSettleRecord, nil)

	entryID, err := l.svcCtx.Scheduler.Register("*/1 * * * *", task)
	if err != nil {
		logx.WithContext(l.ctx).Errorf("!!!MqueueSchedulerErr!!! ====> 【settleRecordScheduler】 registered  err:%+v , task:%+v", err, task)
	}
	fmt.Printf("【settleRecordScheduler】 registered an  entry: %q \n", entryID)
}
