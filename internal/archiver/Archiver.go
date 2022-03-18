package archiver

import (
	"email-archiver-cli/internal/util"
	"email-archiver-cli/proto/emailarchver/email"
	"errors"
	"fmt"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/proto"
)

type Archiver struct {
	configName string
	packager   Packager
	sink       Sink
	log        *log.Entry
}

func NewArchiver(packager Packager, sink Sink, configName string) *Archiver {
	return &Archiver{
		configName: configName,
		packager:   packager,
		sink:       sink,
		log:        util.PrepareLogging(log.Fields{"source": "Archiver", "configName": configName}),
	}
}

func (a *Archiver) NewMail(mail email.EMail) error {
	if validationError := validate(&mail); validationError != nil {
		a.log.Error(validationError)
		return validationError
	}
	a.log.WithFields(log.Fields{"msgId": mail.MessageId}).Info("new mail archiving requested")
	if packaged, err := a.packager.Package(mail); err == nil {
		filename := generateFileName(a.configName, &mail)
		content, _ := proto.Marshal(&packaged)
		a.sink.Persist(content, filename)
	}
	return nil
}

func validate(mail *email.EMail) (err error) {
	if mail == nil {
		err = errors.New("attempted to archive invalid email")
	}
	return
}

func generateFileName(configName string, e *email.EMail) string {
	d := e.Date.AsTime()
	return fmt.Sprintf("%s/%02d/%02d/%02d/%s", configName, d.Year(), d.Month(), d.Day(), util.SanitizeString(e.MessageId))
}
