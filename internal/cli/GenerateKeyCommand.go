package cli

import (
	keygen2 "email-archiver-cli/internal/keygen"
	"email-archiver-cli/internal/repository"
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

type GenerateKeyCommand struct {
	KeyName string
	rotate  bool
}

func (sv *GenerateKeyCommand) Run(ctx *kong.Context) error {
	log.WithFields(log.Fields{
		"keyName": sv.KeyName,
		"rotate":  sv.rotate,
	}).Info("Starting")

	repo, repoErr := repository.Open()

	if repoErr != nil {
		log.WithFields(log.Fields{
			"keyName": sv.KeyName,
			"rotate":  sv.rotate,
			"error":   repoErr,
		}).Fatal("Unable to open repository.")
		return repoErr
	}

	keyRepo := repository.NewFileKeyRepository(*repo)
	keygen, _ := keygen2.NewKeygen(keyRepo)

	key, keygenError := keygen.CreateKey(sv.KeyName)
	if keygenError != nil {
		log.WithFields(log.Fields{
			"keyName": sv.KeyName,
			"rotate":  sv.rotate,
			"error":   keygenError,
		}).Fatal("Unable to create key.")
		return repoErr
	} else {
		log.WithFields(log.Fields{
			"keyName":  sv.KeyName,
			"rotate":   sv.rotate,
			"revision": key.Revision.String(),
		}).Info("Created new key.")
	}

	errPersist := keyRepo.Persist(sv.KeyName, key)
	if errPersist != nil {
		log.WithFields(log.Fields{
			"keyName": sv.KeyName,
			"rotate":  sv.rotate,
			"error":   errPersist,
		}).Fatal("Unable to persist key.")
		return repoErr
	}

	return nil
}
