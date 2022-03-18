package cli

import (
	"email-archiver-cli/internal/keyService"
	"email-archiver-cli/internal/repository"
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

type GenerateKeyCommand struct {
	KeyName string
	Rotate  bool
}

func (sv *GenerateKeyCommand) Run(ctx *kong.Context) error {
	sv.log(log.DebugLevel, "Generating new Key")

	if keyService, err := BuildKeyService(); err != nil {
		sv.log(log.FatalLevel, "Unable to open key repository.")
		return err
	} else {
		if key, keygenError := keyService.CreateKey(sv.KeyName, sv.Rotate); keygenError != nil {
			sv.log(log.FatalLevel, "Unable to create key.")
			return keygenError
		} else {
			sv.log(log.InfoLevel, "New key created", key.Revision.String())
		}
	}

	return nil
}

func (sv *GenerateKeyCommand) log(level log.Level, message ...string) {
	log.WithFields(log.Fields{
		"keyName": sv.KeyName,
		"Rotate":  sv.Rotate,
	}).Log(level, message)
}

func BuildKeyService() (KeyService, error) {
	repo, _ := repository.Open()
	keyRepo := repository.NewFileKeyRepository(*repo)
	keygen, err := keyService.NewKeyService(keyRepo)

	return keygen, err
}
