package keygen

import (
	"crypto/rand"
	"crypto/rsa"
	"email-archiver-cli/internal/repository"
	"email-archiver-cli/models"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type Keygen struct {
	keyRepository repository.KeyRepository
}

const DefaultKeylength = 4096

func (k *Keygen) CreateKey(keyName string) (key *models.Key, err error) {
	log.WithFields(log.Fields{
		"keyName": keyName,
	}).Debug("Generating key")

	if k.keyRepository.Contains(keyName) {
		return nil, KeyAlreadyExists
	}

	if key, err = k.generateKey(DefaultKeylength); err == nil {
		k.keyRepository.Persist(keyName, key)
	}

	return
}

func (k *Keygen) generateKey(keylength int) (key *models.Key, err error) {

	privateKey, err := rsa.GenerateKey(rand.Reader, keylength)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey

	revision, errUuid := uuid.NewRandom()
	if errUuid != nil {
		return nil, err
	}

	key = &models.Key{
		Revision:   revision,
		PublicKey:  publicKey,
		PrivateKey: privateKey,
	}

	return
}

func NewKeygen(repository repository.KeyRepository) (keygen Keygen, err error) {
	return Keygen{keyRepository: repository}, nil
}

var KeyAlreadyExists = errors.New("Key already exists")
