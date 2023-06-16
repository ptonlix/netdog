package timeplugin

import (
	"context"
	"fmt"
	"time"

	"github.com/ptonlix/netdog/internal/pingtest"
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
	ctx    context.Context
}

func NewTimePlugin(spec string, plugin Pluginer, logger *zap.Logger) *TimePlugin {
	nyc, _ := time.LoadLocation("Asia/Shanghai")

	return &TimePlugin{P: plugin, C: cron.New(cron.WithLocation(nyc), cron.WithSeconds()), J: &Job{}, Spec: spec, logger: logger}
}

func (t *TimePlugin) StartCronJob(ctx context.Context) {

	t.ctx = ctx
	fmt.Println(t)
	t.C.AddJob(t.Spec, t)

	// t.C.AddFunc("*/5 * * * * *", func() {

	// 	fmt.Println("每分钟执行一次")
	// })

	// 启动执行任务
	t.C.Start()
	// 退出时关闭计划任务
	defer t.C.Stop()

	<-ctx.Done()
}

func (t *TimePlugin) Run() {
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
		fmt.Printf("12355")
		// go func() {
		// 	<-timer.C
		// 	// 开启Ping探测
		// 	ctx, _ := context.WithTimeout(t.ctx, pt.Durtime)
		// 	test := pingtest.NewPingTestServer(t.logger)
		// 	t.logger.Debug("", zap.String("debug", fmt.Sprintf("%+v", test.LoopTest(ctx))))
		// 	// 发送至记录输出模块
		// }()
		<-timer.C

		// 开启Ping探测
		ctx, _ := context.WithTimeout(t.ctx, pt.Durtime)
		test := pingtest.NewPingTestServer(t.logger)
		t.logger.Debug("", zap.String("debug", fmt.Sprintf("%+v", test.LoopTest(ctx))))
	}

	// bindwidthTimerList := []*time.Timer{}
	// for _, pt := range t.BindwidthTime {
	// 	timer := time.NewTimer(pt.Start.Sub(now))
	// 	bindwidthTimerList = append(bindwidthTimerList, timer)
	// }

}
