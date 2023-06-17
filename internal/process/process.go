package process

import (
	"fmt"
	"os"
	"reflect"
	"time"

	"github.com/ptonlix/netdog/configs"
	"github.com/ptonlix/netdog/internal/pingtest"
	"go.uber.org/zap"
)

//数据处理
type Process struct {
	RecordStart   time.Time
	RecordDurtime time.Duration
	NetworkNode   string
	logger        *zap.Logger
}

func NewProcess(start time.Time, dur time.Duration, logger *zap.Logger) *Process {
	return &Process{RecordStart: start, RecordDurtime: dur, NetworkNode: configs.Get().Network.Node, logger: logger}
}

func (p *Process) WritePingData(data []pingtest.TestDeviceResult) error {
	result := fmt.Sprintf("%s \nPingTest Start time: %+v \nDuration time: %+v\n", p.NetworkNode, p.RecordStart, p.RecordDurtime)
	for _, d := range data {
		typeOfData := reflect.TypeOf(d)
		ValueOfData := reflect.ValueOf(d)
		str := ""
		for i := 0; i < typeOfData.NumField(); i++ {
			// 获取每个成员的结构体字段类型
			fieldType := typeOfData.Field(i)
			str += fmt.Sprintf("%s : %v ", fieldType.Name, ValueOfData.FieldByName(fieldType.Name))
		}
		str += "\n"
		result += str
	}
	// fmt.Println(result)
	//写入文件
	// _, flag := file.IsExists(configs.Get().Data.Pingfile)
	// fmt.Println(configs.Get().Data.Pingfile)
	// if !flag {
	// 	p.logger.Error("File is not Exist", zap.String("filepath", configs.Get().Data.Pingfile))
	// 	return errors.New("file is not exist")
	// }
	err := p.writeFile(configs.Get().Data.Pingfile, result)
	if err != nil {
		return err
	}

	return nil
}

func (p *Process) writeFile(filepath string, data string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_APPEND, 0)
	if err != nil {
		p.logger.Error("OpenFile failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}
	defer file.Close()

	counts, err := file.WriteString(data)
	if err != nil {
		p.logger.Error("Write FileData failed", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}
	p.logger.Info("Write FileData", zap.String("counts", fmt.Sprintf("%+v", counts)))

	return nil
}
