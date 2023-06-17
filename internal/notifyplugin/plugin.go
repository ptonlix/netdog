package notifyplugin

import (
	"fmt"

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

func (n *NotifyPlugin) Nofity(content string) error {
	if err := n.P.Send(content); err != nil {
		n.logger.Error("Nofity test result", zap.String("error", fmt.Sprintf("%+v", err)))
		return err
	}
	return nil
}
