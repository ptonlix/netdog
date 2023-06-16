//go:build darwin

package pingtest

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/ptonlix/netdog/configs"
	"go.uber.org/zap"
)

type TestDeviceResult struct {
	DeviceName  string
	DeviceIp    string
	SendPackets int
	RecvPackets int
	LossPackets float64
	AvgRtt      time.Duration
}

type TestDeviceInfo struct {
	Name string
	Ip   string
}

type PingTestServer struct {
	wg      sync.WaitGroup
	Devices []TestDeviceInfo
	logger  *zap.Logger
}

func NewPingTestServer(logger *zap.Logger) *PingTestServer {
	devs := make([]TestDeviceInfo, len(configs.Get().Network.PingTest.Device))
	for i, d := range configs.Get().Network.PingTest.Device {
		devs[i].Ip = d.Ip
		devs[i].Name = d.Name
	}
	return &PingTestServer{logger: logger, Devices: devs}
}

func (t *PingTestServer) LoopTest(ctx context.Context) []TestDeviceResult {
	var result []TestDeviceResult

	// 函数内的局部变量channel, 专门用来接收函数内所有goroutine的结果
	resultChannel := make(chan TestDeviceResult, 20)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}

	go t.PingHandle(&result, resultChannel, wgResponse)

	for _, devinfo := range t.Devices {
		t.wg.Add(1)
		go t.PingG(ctx, devinfo.Ip, devinfo.Name, resultChannel)
	}
	t.wg.Wait()
	close(resultChannel)
	// 等待wgResponse的计数器归零
	wgResponse.Wait()

	// 返回聚合后结果
	return result
}
func (t *PingTestServer) PingHandle(result *[]TestDeviceResult, resultChannel <-chan TestDeviceResult, s *sync.WaitGroup) {
	// wgResponse计数器+1
	s.Add(1)
	// 读取结果
	for response := range resultChannel {
		// 处理结果
		*result = append(*result, response)
	}
	s.Done()
}

func (t *PingTestServer) PingG(ctx context.Context, ip string, name string, resultChannel chan<- TestDeviceResult) {
	defer t.wg.Done()

	pinger, err := ping.NewPinger(ip)
	if err != nil {
		t.logger.Error("Network Test Error:", zap.String("error", fmt.Sprintf("%+v", err)))
		return
	}
	pinger.SetPrivileged(false)
	go pinger.Run() // Blocks until finished.
	// if err != nil {
	// 	t.logger.Error("Network Test Error:", zap.String("error", fmt.Sprintf("%+v", err)))
	// 	retrn
	// }

	//等待取消信号
	<-ctx.Done()
	pinger.Stop()

	//汇总数据
	result := TestDeviceResult{
		DeviceName: name, DeviceIp: ip,
		SendPackets: pinger.Statistics().PacketsSent,
		RecvPackets: pinger.Statistics().PacketsRecv,
		LossPackets: pinger.Statistics().PacketLoss,
		AvgRtt:      pinger.Statistics().AvgRtt,
	}
	resultChannel <- result
}
