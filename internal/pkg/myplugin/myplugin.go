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
	times, _ := time.Parse(timeLayout, "2023-06-15 18:22:18")

	t := []timeplugin.TestTime{
		{Start: times, Durtime: time.Second * 10},
	}
	return t, nil
}

func (m *MyPlugin) BindwidthTestTime() ([]timeplugin.TestTime, error) {
	// // go语言固定日期模版
	// timeLayout := "2006-01-02 15:04:05"
	// times, _ := time.Parse(timeLayout, "2023-06-15 17:50:18")

	// t := []timeplugin.TestTime{
	// 	{Start: times, Durtime: time.Second * 10},
	// }
	return nil, nil
}
