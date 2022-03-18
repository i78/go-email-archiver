package repository

import (
	"email-archiver-cli/internal/util"
	log "github.com/sirupsen/logrus"
)

type MailRepository struct {
	filesystem Filesystem
	log        *log.Entry
}

func NewMailRepository(parentRepository Repository) *MailRepository {
	return &MailRepository{
		filesystem: parentRepository.filesystem,
		log:        util.PrepareLogging(log.Fields{"source": "FileKeyRepository"})}
}

func (r MailRepository) Persist(content []byte, filenameHint string) error {
	return r.filesystem.Save(r.createFilename(filenameHint), content)
}

func (r MailRepository) createFilename(filenameHint string) string {
	root := createFilename(mailSubfolderPrefix)
	return root + "/" + filenameHint
}
