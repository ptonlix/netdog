package flightstime

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ptonlix/netdog/internal/timeplugin"
	"github.com/ptonlix/netdog/pkg/httpclient"
	"go.uber.org/zap"

	httpURL "net/url"
)

const version = 2

type Flightsinfo struct {
	LineName   string    //航线
	DepartTime time.Time //开航时间
	ArriveTime time.Time //达到时间
}

type FlightsMap map[string][]Flightsinfo

type Flightstime struct {
	Url     string
	Flights FlightsMap
	logger  *zap.Logger
}

func NewFlightstime(url string, logger *zap.Logger) *Flightstime {
	return &Flightstime{Url: url, logger: logger, Flights: make(FlightsMap)}
}
func (f *Flightstime) getTodayFlightTime() error {
	form := httpURL.Values{}
	form.Set("version", "2")
	body, err := httpclient.Get(f.Url, form)
	if err != nil {
		f.logger.Error("Get TodayFlightTime Failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}
	//解析 body
	todaydata := Todyflightsdata{}

	if err := json.Unmarshal(body, &todaydata); err != nil {
		f.logger.Error("Unmarshal TodayFlightTime Failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}

	// 处理北涠航线
	for _, v := range todaydata.Data.BW {
		timelayout := "2006-01-02 15:04"
		Departtime, _ := time.Parse(timelayout, v.Date+" "+v.DepartTime)
		Arrivetime, _ := time.Parse(timelayout, v.Date+" "+v.ArriveTime)
		f.Flights[string(v.ShipName)] = append(f.Flights[string(v.ShipName)], Flightsinfo{LineName: string(v.LineName), DepartTime: Departtime, ArriveTime: Arrivetime})
	}
	// 处理涠北航线
	for _, v := range todaydata.Data.Wb {
		timelayout := "2006-01-02 15:04"
		Departtime, _ := time.Parse(timelayout, v.Date+" "+v.DepartTime)
		Arrivetime, _ := time.Parse(timelayout, v.Date+" "+v.ArriveTime)
		f.Flights[string(v.ShipName)] = append(f.Flights[string(v.ShipName)], Flightsinfo{LineName: string(v.LineName), DepartTime: Departtime, ArriveTime: Arrivetime})
	}
	return nil
}

func (f *Flightstime) PingTestTime() ([]timeplugin.TestTime, error) {

	return nil, nil
}
func (f *Flightstime) BindwidthTestTime() ([]timeplugin.TestTime, error) {
	return nil, nil
}
