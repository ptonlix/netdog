package dogmail

import (
	"fmt"

	"github.com/ptonlix/netdog/configs"
	"github.com/ptonlix/netdog/pkg/mail"
)

type Dogmail struct {
	o mail.Options
}

func NewDogmail() *Dogmail {

	dogmail := Dogmail{}
	dogmail.o.MailHost = configs.Get().Mail.Host
	dogmail.o.MailPort = configs.Get().Mail.Port
	dogmail.o.MailUser = configs.Get().Mail.User
	dogmail.o.MailPass = configs.Get().Mail.Pass
	dogmail.o.MailTo = configs.Get().Mail.To
	dogmail.o.Subject = fmt.Sprintf("%s Network KPI Data", configs.Get().Network.Node)

	return &dogmail
}

func (d *Dogmail) Send(content string) error {
	d.o.Body = content
	if err := mail.Send(&d.o); err != nil {
		return err
	}
	return nil
}
