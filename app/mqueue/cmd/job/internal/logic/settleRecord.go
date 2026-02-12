package logic

import (
	"context"
	"fmt"

	"github.com/hibiken/asynq"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/internal/svc"
)

type SettleRecordHandler struct {
	svcCtx *svc.ServiceContext
}

func NewSettleRecordHandler(svcCtx *svc.ServiceContext) *SettleRecordHandler {
	return &SettleRecordHandler{
		svcCtx: svcCtx,
	}
}

// ---------------------------
// @brief 测试任务,每一分钟执行一次
// ---------------------------
func (l *SettleRecordHandler) ProcessTask(ctx context.Context, _ *asynq.Task) error {

	fmt.Printf("shcedule job demo -----> every one minute exec \n")

	return nil
}
