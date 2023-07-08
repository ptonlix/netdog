package flightstime

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ptonlix/netdog/configs"
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
	f := Flightstime{Url: url, logger: logger, Flights: make(FlightsMap)}
	return &f
}
func (f *Flightstime) GetTodayFlightTime() error {
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
	loc, err := time.LoadLocation("Local")
	if err != nil {
		f.logger.Error("SelectTimeLocation Failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}
	// 处理北涠航线
	for _, v := range todaydata.Data.BW {
		timelayout := "2006-01-02 15:04"
		Departtime, _ := time.ParseInLocation(timelayout, v.Date+" "+v.DepartTime, loc)
		Arrivetime, _ := time.ParseInLocation(timelayout, v.Date+" "+v.ArriveTime, loc)
		f.Flights[string(v.ShipName)] = append(f.Flights[string(v.ShipName)], Flightsinfo{LineName: string(v.LineName), DepartTime: Departtime, ArriveTime: Arrivetime})
	}
	// 处理涠北航线
	for _, v := range todaydata.Data.Wb {
		timelayout := "2006-01-02 15:04"
		Departtime, _ := time.ParseInLocation(timelayout, v.Date+" "+v.DepartTime, loc)
		Arrivetime, _ := time.ParseInLocation(timelayout, v.Date+" "+v.ArriveTime, loc)
		f.Flights[string(v.ShipName)] = append(f.Flights[string(v.ShipName)], Flightsinfo{LineName: string(v.LineName), DepartTime: Departtime, ArriveTime: Arrivetime})
	}
	return nil
}

func (f *Flightstime) PingTestTime() ([]timeplugin.TestTime, error) {
	// 获取获取航班时间
	if len(f.Flights) == 0 {
		if err := f.GetTodayFlightTime(); err != nil {
			f.logger.Error("PingTest Time Error", zap.String("error", fmt.Sprintf("%+v", err)))
			return nil, err
		}
	}
	if _, ok := f.Flights[configs.Get().Network.Node]; !ok {
		return nil, fmt.Errorf("%s is not running today", configs.Get().Network.Node)
	}
	infoList := f.Flights[configs.Get().Network.Node]
	f.logger.Info("Flightstime InfoList:", zap.String("list", fmt.Sprintf("%+v", infoList)))
	testtime := make([]timeplugin.TestTime, len(infoList))
	for i := 0; i < len(infoList); i++ {
		testtime[i].Start = infoList[i].DepartTime
		testtime[i].Durtime = infoList[i].ArriveTime.Sub(infoList[i].DepartTime)
		testtime[i].Detail = infoList[i].LineName
	}
	return testtime, nil
}
func (f *Flightstime) BindwidthTestTime() ([]timeplugin.TestTime, error) {
	// 获取获取航班时间
	if len(f.Flights) == 0 {
		if err := f.GetTodayFlightTime(); err != nil {
			f.logger.Error("PingTest Time Error", zap.String("error", fmt.Sprintf("%+v", err)))
			return nil, err
		}
	}
	if _, ok := f.Flights[configs.Get().Network.Node]; !ok {
		return nil, fmt.Errorf("%s is not running today", configs.Get().Network.Node)
	}
	infoList := f.Flights[configs.Get().Network.Node]
	testtime := make([]timeplugin.TestTime, len(infoList)*2)
	strDurtime := configs.Get().Network.BindwidthTest.BindwidthDuration
	durtime, err := time.ParseDuration(strDurtime)
	if err != nil {
		f.logger.Error("Parse Duration Time Error", zap.String("error", fmt.Sprintf("%+v", err)))
		return nil, err
	}
	for i := 0; i < len(infoList); i++ {
		flightsDur := infoList[i].ArriveTime.Sub(infoList[i].DepartTime) / 4

		testtime[2*i].Start = infoList[i].DepartTime.Add(flightsDur)
		testtime[2*i].Durtime = durtime
		testtime[2*i].Detail = infoList[i].LineName

		testtime[2*i+1].Start = infoList[i].DepartTime.Add(3 * flightsDur)
		testtime[2*i+1].Durtime = durtime
		testtime[2*i+1].Detail = infoList[i].LineName
	}
	return testtime, nil
}
