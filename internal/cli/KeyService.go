package cli

import "email-archiver-cli/models"

type KeyService interface {
	Contains(keyName string) bool
	CreateKey(keyName string, rotate bool) (key *models.Key, err error)
}
