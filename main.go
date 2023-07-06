package main

import (
	"context"
	"fmt"

	"github.com/ptonlix/netdog/configs"
	"github.com/ptonlix/netdog/internal/pingtest"
	"github.com/ptonlix/netdog/internal/pkg/flightstime"
	"github.com/ptonlix/netdog/internal/timeplugin"
	"github.com/ptonlix/netdog/pkg/env"
	"github.com/ptonlix/netdog/pkg/logger"
	"github.com/ptonlix/netdog/pkg/shutdown"
	"github.com/ptonlix/netdog/pkg/timeutil"
)

func main() {
	// 初始化 access logger
	accessLogger, err := logger.NewJSONLogger(
		logger.WithDisableConsole(),
		logger.WithField("domain", fmt.Sprintf("%s[%s]", configs.ProjectName, env.Active().Value())),
		logger.WithTimeLayout(timeutil.CSTLayout),
		logger.WithFileP(configs.ProjectLogFile),
	)
	if err != nil {
		panic(err)
	}
	dev := make([]pingtest.TestDeviceInfo, len(configs.Get().Network.PingTest.Device))
	for i, d := range configs.Get().Network.PingTest.Device {
		dev[i].Ip = d.Ip
		dev[i].Name = d.Name
	}
	ctx := context.Background()
	apihost := configs.Get().Network.Api
	ft := flightstime.NewFlightstime(apihost, accessLogger)
	cronTask := timeplugin.NewTimePlugin(configs.Get().Network.Cron, ft, accessLogger)
	accessLogger.Info("The NetDog Start...")
	cronTask.StartCronJob(ctx)

	defer func() {
		_ = accessLogger.Sync()
	}()
	// 优雅关闭
	shutdown.NewHook().Close(
		func() {
			accessLogger.Info("Close The NetDog...")
			ctx.Done()
		},
	)
}
