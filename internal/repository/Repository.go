package repository

import (
	"errors"
	"os"
)

const folderPrefix = "mailArchive"

type MailRepository struct {
}

func InitRepository() (repo *MailRepository, err error) {
	return runIfRepoDoesNotExist(createRepository)
}

func runIfRepoDoesNotExist(fn func() (repo *MailRepository, err error)) (repo *MailRepository, err error) {
	if _, statError := os.Stat(folderPrefix); !os.IsNotExist(statError) {
		err = RepositoryAlreadyExistsError
		return
	}
	return fn()
}

func createRepository() (repo *MailRepository, err error) {
	folders := []string{
		folderPrefix,
		folderPrefix + "/keys",
		folderPrefix + "/maildir",
		folderPrefix + "/index",
	}

	for _, folder := range folders {
		err := os.Mkdir(folder, 0700)
		if err != nil {
			return nil, err
		}
	}

	return
}

var RepositoryAlreadyExistsError = errors.New("Repository already exists. No need to create a new one")
