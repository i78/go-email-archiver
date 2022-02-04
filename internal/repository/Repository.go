package repository

import (
	"errors"
	"os"
)

const folderPrefix = "mailArchive"
const keySubfolderPrefix = "keys"

type MailRepository struct {
}

func Open() (repo *MailRepository, err error) {
	if repoExistsInWorkingDirectory() {
		repo = &MailRepository{}
	} else {
		err = RepositoryNotFoundError
	}
	return
}

func InitRepository() (repo *MailRepository, err error) {
	return runIfRepoDoesNotExist(createRepository)
}

func runIfRepoDoesNotExist(fn func() (repo *MailRepository, err error)) (repo *MailRepository, err error) {
	if repoExistsInWorkingDirectory() {
		err = RepositoryAlreadyExistsError
		return
	}
	return fn()
}

func repoExistsInWorkingDirectory() bool {
	if _, statError := os.Stat(folderPrefix); !os.IsNotExist(statError) {
		return true
	}
	return false
}

func (r *MailRepository) FolderExists(path string) bool {
	if _, statError := os.Stat(path); !os.IsNotExist(statError) {
		return true
	}
	return false
}

func createRepository() (repo *MailRepository, err error) {
	folders := []string{
		folderPrefix,
		folderPrefix + "/" + keySubfolderPrefix,
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
var RepositoryNotFoundError = errors.New("No Repository found. You might want to create one")
