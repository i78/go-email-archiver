package repository

import (
	"email-archiver-cli/models"
)

type KeyRepository interface {
	Contains(keyName string) bool
	Persist(keyName string, key *models.Key) error
}

// var UnableToResolveRepository = errors.New("Unable to resolve repository.")
