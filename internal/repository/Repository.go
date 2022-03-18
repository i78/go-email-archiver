package repository

import (
	"email-archiver-cli/internal/filesystem"
	"email-archiver-cli/internal/util"
	"errors"
	"os"
	"strings"
)

const folderPrefix = "mailArchive"
const keySubfolderPrefix = "keys"
const mailSubfolderPrefix = "maildir"
const configSubfolderPrefix = "config"

type Repository struct {
	filesystem Filesystem
}

func Open() (repo *Repository, err error) {
	if repoExistsInWorkingDirectory() {
		repo = &Repository{
			filesystem: filesystem.LocalFilesystem{},
		}
	} else {
		err = NotFoundError
	}
	return
}

func InitRepository() (repo *Repository, err error) {
	return whenRepositoryDoesNotExist(createRepository)
}

func whenRepositoryDoesNotExist(fn func() (repo *Repository, err error)) (repo *Repository, err error) {
	if repoExistsInWorkingDirectory() {
		err = AlreadyExistsError
		return
	}
	return fn()
}

func repoExistsInWorkingDirectory() bool {
	return filesystem.Exists(folderPrefix)
}

func createFilename(tokens ...string) string {
	sanitized := util.Sanitize(tokens)
	return folderPrefix + "/" + strings.Join(sanitized, "/")
}

func createRepository() (repo *Repository, err error) {
	folders := []string{
		createFilename(),
		createFilename(keySubfolderPrefix),
		createFilename(configSubfolderPrefix),
		createFilename("/maildir"),
		createFilename("/index"),
	}

	util.ForEach(folders, func(f string) {
		err = os.Mkdir(f, 0700)
	})

	return
}

var AlreadyExistsError = errors.New("repository already exists. No need to create a new one")
var NotFoundError = errors.New("no Repository found. You might want to create one")
