package process

import (
	"testing"
	"time"

	"github.com/ptonlix/netdog/internal/pingtest"
	"github.com/ptonlix/netdog/pkg/logger"
)

func TestWritePingData(t *testing.T) {
	accessLogger, _ := logger.NewJSONLogger()
	accessLogger.Info("test")
	data := []pingtest.TestDeviceResult{
		{DeviceName: "test", DeviceIp: "192.168.1.1", SendPackets: 100, RecvPackets: 100, LossPackets: 0, AvgRtt: time.Millisecond * 20},
		{DeviceName: "test", DeviceIp: "192.168.1.1", SendPackets: 100, RecvPackets: 100, LossPackets: 0, AvgRtt: time.Millisecond * 20},
		{DeviceName: "test", DeviceIp: "192.168.1.1", SendPackets: 100, RecvPackets: 100, LossPackets: 0, AvgRtt: time.Millisecond * 20},
	}
	process := NewProcess(time.Now(), time.Minute*60, accessLogger)
	err := process.WritePingData(data)
	if err != nil {
		t.Error("write ping data error", err)
		return
	}
	t.Log("success")
}
