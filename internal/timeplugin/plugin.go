package timeplugin

import (
	"context"
	"fmt"
	"time"

	"github.com/ptonlix/netdog/internal/pingtest"
	"github.com/ptonlix/netdog/internal/process"
	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

type TestTime struct {
	Start   time.Time
	Durtime time.Duration
}
type Pluginer interface {
	PingTestTime() ([]TestTime, error)
	BindwidthTestTime() ([]TestTime, error)
}
type Job struct {
}

func (j *Job) Run() {
	//获取当前测试时间
}

type TimePlugin struct {
	P    Pluginer
	C    *cron.Cron
	J    *Job
	Spec string //定时任务Cron表达式

	PingTime      []TestTime
	BindwidthTime []TestTime

	logger *zap.Logger
}

func NewTimePlugin(spec string, plugin Pluginer, logger *zap.Logger) *TimePlugin {

	return &TimePlugin{
		P: plugin,
		C: cron.New(cron.WithLocation(time.Local), cron.WithSeconds()),
		J: &Job{}, Spec: spec, logger: logger}
}

func (t *TimePlugin) StartCronJob(ctx context.Context) {
	//t.C.AddJob(t.Spec, t)
	t.C.AddFunc(t.Spec, func() {
		t.Run(ctx)
	})
	// 启动执行任务
	t.C.Start()
	// 退出时关闭计划任务
	defer t.C.Stop()

	<-ctx.Done()
}

func (t *TimePlugin) Run(ctx context.Context) {
	var err error
	fmt.Println("Start Run")
	t.PingTime, err = t.P.PingTestTime()
	if err != nil {
		t.logger.Error("Get PingTest Time Error:", zap.String("error", fmt.Sprintf("%+v", err)))
		return
	}
	t.BindwidthTime, err = t.P.BindwidthTestTime()
	if err != nil {
		t.logger.Error("Get Bindwidth Time Error:", zap.String("error", fmt.Sprintf("%+v", err)))
		return
	}

	//启动定时任务
	now := time.Now()
	for _, pt := range t.PingTime {
		timer := time.NewTimer(pt.Start.Sub(now))
		t.logger.Info("Time remaining until the next scheduled task is executed:", zap.String("subtime", fmt.Sprintf("%+v", pt.Start.Sub(now))))
		go func() {
			<-timer.C
			t.logger.Info("The Ping detection task starts", zap.String("nowtime", fmt.Sprintf("%+v", time.Now())))
			// 开启Ping探测
			ctx, _ := context.WithTimeout(ctx, pt.Durtime)
			test := pingtest.NewPingTestServer(t.logger)
			pingresult := test.LoopTest(ctx)
			t.logger.Info("", zap.String("debug", fmt.Sprintf("%+v", pingresult)))
			// 发送至记录输出模块
			pro := process.NewProcess(pt.Start, pt.Durtime, t.logger)
			pro.WritePingData(pingresult)
		}()
	}

	// bindwidthTimerList := []*time.Timer{}
	// for _, pt := range t.BindwidthTime {
	// 	timer := time.NewTimer(pt.Start.Sub(now))
	// 	bindwidthTimerList = append(bindwidthTimerList, timer)
	// }

}
