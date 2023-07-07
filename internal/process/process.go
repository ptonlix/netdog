package process

import (
	"fmt"
	"os"
	"reflect"
	"sync"
	"time"

	"github.com/ptonlix/netdog/configs"
	"github.com/ptonlix/netdog/internal/bindwidthtest"
	"github.com/ptonlix/netdog/internal/pingtest"
	"go.uber.org/zap"
)

//数据处理
type Process struct {
	NetworkNode string
	logger      *zap.Logger
	mutex       *sync.Mutex
}

func NewProcess(logger *zap.Logger) *Process {
	return &Process{NetworkNode: configs.Get().Network.Node, logger: logger,
		mutex: &sync.Mutex{}}
}

func (p *Process) WriteBindwidthData(recordStart time.Time, recordDurtime time.Duration, data []bindwidthtest.TestDeviceResult) error {
	result := fmt.Sprintf("%s \nPingTest Start time: %+v \nDuration time: %+v\n", p.NetworkNode, recordStart, recordDurtime)
	for _, v := range data {
		result += string(v)
	}
	err := p.writeFile(configs.Get().Data.Bindwidthfile, result)
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) WritePingData(recordStart time.Time, recordDurtime time.Duration, data []pingtest.TestDeviceResult) error {
	result := fmt.Sprintf("%s \nPingTest Start time: %+v \nDuration time: %+v\n", p.NetworkNode, recordStart, recordDurtime)

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
	//p.mutex.Lock()
	err := p.writeFile(configs.Get().Data.Pingfile, result)
	if err != nil {
		return err
	}
	//p.mutex.Unlock()
	return nil
}

func (p *Process) writeFile(filepath string, data string) error {
	file, err := os.OpenFile(filepath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
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
	file.Sync()
	p.logger.Info("Write FileData", zap.String("counts", fmt.Sprintf("%+v", counts)))

	return nil
}
