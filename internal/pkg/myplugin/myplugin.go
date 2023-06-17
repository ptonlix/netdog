package myplugin

import (
	"time"

	"github.com/ptonlix/netdog/internal/timeplugin"
)

type MyPlugin struct {
}

func (m *MyPlugin) PingTestTime() ([]timeplugin.TestTime, error) {
	// go语言固定日期模版
	timeLayout := "2006-01-02 15:04:05"
	times, _ := time.ParseInLocation(timeLayout, "2023-06-16 22:30:18", time.Local)

	t := []timeplugin.TestTime{
		{Start: times, Durtime: time.Second * 5},
	}
	return t, nil
}

func (m *MyPlugin) BindwidthTestTime() ([]timeplugin.TestTime, error) {
	// go语言固定日期模版
	timeLayout := "2006-01-02 15:04:05"
	times, _ := time.ParseInLocation(timeLayout, "2023-06-16 22:30:18", time.Local)

	t := []timeplugin.TestTime{
		{Start: times, Durtime: time.Second * 5},
	}
	return t, nil
}
