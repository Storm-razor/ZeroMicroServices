package svc

import (
	"fmt"
	"time"

	"github.com/hibiken/asynq"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/scheduler/internal/config"
)

// ---------------------------
// @brief 创建asynq的schduler
// ---------------------------
func newScheduler(c config.Config) *asynq.Scheduler {
	location, _ := time.LoadLocation("Asia/Shanghai")
	return asynq.NewScheduler(
		asynq.RedisClientOpt{
			Addr:     c.Redis.Host,
			Password: c.Redis.Pass,
		}, &asynq.SchedulerOpts{
			Location: location,
			EnqueueErrorHandler: func(task *asynq.Task, opts []asynq.Option, err error) {
				fmt.Printf("Scheduler EnqueueErrorHandler <<<<<<<===>>>>> err : %+v , task : %+v", err, task)
			},
		},
	)
}
