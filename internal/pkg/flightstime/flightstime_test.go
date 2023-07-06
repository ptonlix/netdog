package flightstime

import (
	"fmt"
	"testing"

	"github.com/ptonlix/netdog/pkg/logger"
)

func TestGetTodayFlightTime(t *testing.T) {
	accessLogger, _ := logger.NewJSONLogger()
	f := NewFlightstime("https://m.laiu8.cn/intelligentapi/todayflightsdata", accessLogger)
	err := f.GetTodayFlightTime()
	if err != nil {
		t.Error("Get TodayFilghtTime Failed!", err)
		return
	}
	fmt.Println(f.Flights)
	t.Log("Successfully")
}

func TestPingTestTime(t *testing.T) {
	accessLogger, _ := logger.NewJSONLogger()
	f := NewFlightstime("https://m.laiu8.cn/intelligentapi/todayflightsdata", accessLogger)
	err := f.GetTodayFlightTime()
	if err != nil {
		t.Error("Get TodayFlightTime Failed!", err)
		return
	}
	if testlist, err := f.PingTestTime(); err != nil {
		t.Error("Get PingTestTime Failed!", err)
		return
	} else {
		fmt.Println(testlist)
	}
	t.Log("Successfully")
}
