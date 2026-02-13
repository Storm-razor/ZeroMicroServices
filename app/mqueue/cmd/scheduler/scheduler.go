package main

import (
	"context"
	"flag"
	"os"

	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/scheduler/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/scheduler/internal/logic"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/scheduler/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/scheduler.yaml", "Specify the config file")

func main() {
	flag.Parse()
	var c config.Config

	conf.MustLoad(*configFile, &c)

	logx.DisableStat()

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()

	mqueueScheduler := logic.NewCronScheduler(ctx, svcCtx)
	mqueueScheduler.Register()

	if err := svcCtx.Scheduler.Run(); err != nil {
		logx.Errorf("!!!MqueueSchedulerErr!!!  run err:%+v", err)
		os.Exit(1)
	}
}
