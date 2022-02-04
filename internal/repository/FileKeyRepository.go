package repository

import (
	"crypto/x509"
	"email-archiver-cli/models"
	"io/ioutil"
	"os"
)

type FileKeyRepository struct {
	parentRepository MailRepository
}

func NewFileKeyRepository(parentRepository MailRepository) *FileKeyRepository {
	return &FileKeyRepository{parentRepository: parentRepository}
}

func (f *FileKeyRepository) Contains(keyName string) bool {
	folderName := folderPrefix + "/" + keySubfolderPrefix + "/" + keyName // sanitize todo
	return f.parentRepository.FolderExists(folderName)
}

func (f *FileKeyRepository) Persist(keyName string, key *models.Key) error {
	folderName := folderPrefix + "/" + keySubfolderPrefix + "/" + keyName // sanitize todo
	privateKeyFileName := folderName + "/" + key.Revision.String() + ".key"
	publicKeyFileName := folderName + "/" + key.Revision.String() + ".pub"

	os.Mkdir(folderName, 0700)
	pub := x509.MarshalPKCS1PublicKey(key.PublicKey)
	errPub := ioutil.WriteFile(publicKeyFileName, pub, 0644)
	if errPub != nil {
		return errPub
	}

	prv := x509.MarshalPKCS1PrivateKey(key.PrivateKey)
	errPk := ioutil.WriteFile(privateKeyFileName, prv, 0644)
	if errPk != nil {
		return errPk
	}

	return nil
}
