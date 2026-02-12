package main

import (
	"context"
	"flag"
	"os"

	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/internal/config"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/internal/logic"
	"github.com/wwwzy/ZeroMicroServices/app/mqueue/cmd/job/internal/svc"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/logx"
)

var configFile = flag.String("f", "etc/mqueue.yaml", "Specify the config file")

func main() {
	flag.Parse()

	var c config.Config

	conf.MustLoad(*configFile, &c, conf.UseEnv())

	if err := c.SetUp(); err != nil {
		panic(err)
	}

	svcContext := svc.NewServiceContext(c)
	ctx := context.Background()
	cronJob := logic.NewCronJob(ctx, svcContext)
	mux := cronJob.Register()

	if err := svcContext.AsynqServer.Run(mux); err != nil {
		logx.WithContext(ctx).Errorf("!!!CronJobErr!!! run err:%+v", err)
		os.Exit(1)
	}
}
