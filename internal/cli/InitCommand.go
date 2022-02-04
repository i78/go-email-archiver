package cli

import (
	"email-archiver-cli/internal/repository"
	"errors"
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

type InitCommand struct {
}

func (sv *InitCommand) Run(ctx *kong.Context) error {
	log.WithFields(log.Fields{}).Debug("Trying to create new Repository")

	if _, err := repository.InitRepository(); err != nil {
		if errors.Is(err, repository.RepositoryAlreadyExistsError) {
			log.Error("The repository already exists. No need to init it again, just start working with it.")
		} else {
			log.WithFields(log.Fields{"error": err}).Error("Unexpected Error occourred")
		}
		return err
	} else {
		log.WithFields(log.Fields{}).Info("Repository created successfully.")
	}

	return nil
}
