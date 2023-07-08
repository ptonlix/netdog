package notifyplugin

import (
	"fmt"
	"io/ioutil"

	"github.com/ptonlix/netdog/configs"
	"go.uber.org/zap"
)

type Pluginer interface {
	Send(content string) error
}

type NotifyPlugin struct {
	P      Pluginer
	logger *zap.Logger
}

func NewNotifyPlugin(p Pluginer, logger *zap.Logger) *NotifyPlugin {
	return &NotifyPlugin{P: p, logger: logger}
}

func (n *NotifyPlugin) NofityFromDatafile() error {
	//读取ping的记录文件
	pingdata, err := n.readDataFromFile(configs.Get().Data.Pingfile)
	if err != nil {
		n.logger.Warn("read data from file failed", zap.String("error", fmt.Sprintf("%+v", err)))
	}
	bindwidthdata, err := n.readDataFromFile(configs.Get().Data.Bindwidthfile)
	if err != nil {
		n.logger.Warn("read data from file failed", zap.String("error", fmt.Sprintf("%+v", err)))
	}
	testResult := string(pingdata) + string(bindwidthdata)
	if testResult == "" {
		return fmt.Errorf("data empty, nofity failed")
	}
	if err := n.Nofity(testResult); err != nil {
		return err
	}
	return nil
}

func (n *NotifyPlugin) readDataFromFile(filepath string) ([]byte, error) {
	//读取文件
	bytes, err := ioutil.ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

func (n *NotifyPlugin) Nofity(content string) error {
	if err := n.P.Send(content); err != nil {
		n.logger.Error("Nofity test result", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}
	return nil
}
