package flightstime

// 今日航班列表对象
type Todyflightsdata struct {
	Code int64 `json:"code"`
	Data Data  `json:"data"`
}

type Data struct {
	Today   string    `json:"today"`
	BW      []BW      `json:"bw"`
	Wb      []BW      `json:"wb"`
	Nonsail []Nonsail `json:"nonsail"`
	Test    []Test    `json:"test"`
}

type BW struct {
	ID              int64          `json:"id"`
	LineID          int64          `json:"lineID"`
	LineName        LineName       `json:"lineName"`
	Date            string         `json:"date"`
	Code            string         `json:"code"`
	ShipName        Test           `json:"shipName"`
	DepartTime      string         `json:"departTime"`
	ArriveTime      string         `json:"arriveTime"`
	Status          Status         `json:"status"`
	OldShip         string         `json:"oldShip"`
	NewShip         string         `json:"newShip"`
	RealDepartTime  string         `json:"realDepartTime"`
	RealArriveTime  RealArriveTime `json:"realArriveTime"`
	Info            string         `json:"info"`
	DepartCity      City           `json:"departCity"`
	ArriveCity      City           `json:"arriveCity"`
	IsDelay         int64          `json:"isDelay"`
	Delay           int64          `json:"Delay"`
	ArrDelay        int64          `json:"ArrDelay"`
	DepartTimeDelay string         `json:"departTimeDelay"`
	ArriveTimeDelay string         `json:"arriveTimeDelay"`
	ChangeShip      int64          `json:"changeShip"`
	RealShip        Test           `json:"realShip"`
	CxMmsi          *int64         `json:"cx_mmsi,omitempty"`
	CxShipname      *string        `json:"cx_shipname,omitempty"`
	CxDest          *string        `json:"cx_dest,omitempty"`
	CxHdg           *int64         `json:"cx_hdg,omitempty"`
	IsDepart        int64          `json:"isDepart"`
	IsArrive        int64          `json:"isArrive"`
	Lat             int64          `json:"lat"`
	Lon             int64          `json:"lon"`
	Speed           int64          `json:"speed"`
	Progress        int64          `json:"progress"`
	DepartLat       int64          `json:"departLat"`
	ArriveLat       int64          `json:"arriveLat"`
	ShowLoction     int64          `json:"showLoction"`
	OldShipName     string         `json:"old_ship_name"`
	IsMy            int64          `json:"is_my"`
}

type Nonsail struct {
	Shipname   string  `json:"shipname"`
	Lat        *int64  `json:"lat"`
	Lon        *int64  `json:"lon"`
	CxMmsi     *int64  `json:"cx_mmsi"`
	CxShipname *string `json:"cx_shipname"`
	CxDest     *string `json:"cx_dest"`
	CxHdg      *int64  `json:"cx_hdg"`
}

type City string

const (
	北海 City = "北海"
	涠洲 City = "涠洲"
)

type LineName string

const (
	北海涠洲 LineName = "北海-涠洲"
	涠洲北海 LineName = "涠洲-北海"
)

type RealArriveTime string

const (
	Empty RealArriveTime = "- -"
)

type Test string

const (
	北游12 Test = "北游12"
	北游25 Test = "北游25"
	北游26 Test = "北游26"
)

type Status string

const (
	已出发 Status = "已出发"
	未开航 Status = "未开航"
)
