package keyService

import (
	"email-archiver-cli/models"
)

type KeyRepository interface {
	Contains(keyName string) bool
	Persist(keyName string, key *models.Key) error
}
