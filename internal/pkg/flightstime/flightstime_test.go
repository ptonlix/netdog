package flightstime

import (
	"testing"

	"github.com/ptonlix/netdog/pkg/logger"
)

func TestGetTodayFlightTime(t *testing.T) {
	accessLogger, _ := logger.NewJSONLogger()
	f := NewFlightstime("https://m.laiu8.cn/intelligentapi/todayflightsdata", accessLogger)
	if err := f.getTodayFlightTime(); err != nil {
		t.Error("Get TodayFilghtTime Failed!")
	}
	t.Log("Successfully")
}
