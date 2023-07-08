package bindwidthtest

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"sync"

	"github.com/ptonlix/netdog/configs"
	"go.uber.org/zap"
)

type TestDeviceResult string

type TestDeviceInfo struct {
	Name string
	Ip   string
}

type BindwidthTestServer struct {
	wg          sync.WaitGroup
	Devices     []TestDeviceInfo
	logger      *zap.Logger
	Tool        string
	TestDurtime string
}

func NewBindwidthTestServer(logger *zap.Logger) *BindwidthTestServer {
	devs := make([]TestDeviceInfo, len(configs.Get().Network.BindwidthTest.Device))
	for i, d := range configs.Get().Network.BindwidthTest.Device {
		devs[i].Ip = d.Ip
		devs[i].Name = d.Name
	}
	tool := configs.Get().Network.BindwidthTest.Testtool
	durtime := configs.Get().Network.BindwidthTest.BindwidthDuration
	return &BindwidthTestServer{logger: logger, Devices: devs, Tool: tool, TestDurtime: durtime}
}
func (b *BindwidthTestServer) LoopTest(ctx context.Context) []TestDeviceResult {
	var result []TestDeviceResult

	// 函数内的局部变量channel, 专门用来接收函数内所有goroutine的结果
	resultChannel := make(chan TestDeviceResult, 20)
	// 为读取结果控制器创建新的WaitGroup, 需要保证控制器内的所有值都已经正确处理完毕, 才能结束
	wgResponse := &sync.WaitGroup{}

	go b.BindwidthHandle(&result, resultChannel, wgResponse)

	for _, devinfo := range b.Devices {
		b.wg.Add(1)
		go b.BindwidthTestG(ctx, devinfo.Ip, devinfo.Name, resultChannel)
	}
	b.wg.Wait()
	close(resultChannel)
	// 等待wgResponse的计数器归零
	wgResponse.Wait()

	// 返回聚合后结果
	return result
}

func (t *BindwidthTestServer) BindwidthHandle(result *[]TestDeviceResult, resultChannel <-chan TestDeviceResult, s *sync.WaitGroup) {
	// wgResponse计数器+1
	s.Add(1)
	// 读取结果
	for response := range resultChannel {
		// 处理结果
		*result = append(*result, response)
	}
	s.Done()
}

func (b *BindwidthTestServer) BindwidthTestG(ctx context.Context, ip string, name string, resultChannel chan<- TestDeviceResult) {
	defer b.wg.Done()
	// 运行ethr命令
	cmd := exec.CommandContext(ctx, b.Tool, "-c", ip, "-n", "0", "-o", configs.EthrLogFile)
	//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		b.logger.Error("Get Cmd stdoutPipe Failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return
	}

	// 启动ethr进程
	err = cmd.Start()
	if err != nil {
		b.logger.Error("Run Cmd Start Failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return
	}

	// 读取并解析ethr输出
	bindwidthContent := ""
	scanner := bufio.NewScanner(stdout)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "SUM") {

			bindwidthContent += line + "\n"
		}
	}

	// 检查是否有读取输出时的错误
	if err := scanner.Err(); err != nil {
		b.logger.Error("Read Cmd Output Error", zap.String("error", fmt.Sprintf("%+v", err)))
	}

	// 等待ethr进程退出
	err = cmd.Wait()
	if err != nil {
		b.logger.Error("Exit Cmd Error", zap.String("error", fmt.Sprintf("%+v", err)))
	}

	bindwidthContent = "TestDevice: " + name + "\n" + bindwidthContent
	resultChannel <- TestDeviceResult(bindwidthContent)
}
