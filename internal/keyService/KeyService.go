package keyService

import (
	"crypto/rand"
	"crypto/rsa"
	"email-archiver-cli/models"
	"errors"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
)

type KeyService struct {
	keyRepository KeyRepository
	keyGenerator  func(int) (key *models.Key, err error)
}

const DefaultKeylength = 4096

func (k KeyService) Contains(keyName string) bool {
	return k.keyRepository.Contains(keyName)
}

func (k KeyService) CreateKey(keyName string, rotate bool) (key *models.Key, err error) {
	log.WithFields(log.Fields{
		"keyName": keyName,
	}).Debug("Generating key")

	if !rotate && k.keyRepository.Contains(keyName) {
		return nil, KeyAlreadyExists
	}

	if key, err = k.keyGenerator(DefaultKeylength); err == nil {
		k.keyRepository.Persist(keyName, key)
	}

	return
}

func genRsaKey(keysize int) (key *models.Key, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, keysize)
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

func NewKeyService(repository KeyRepository) (keygen KeyService, err error) {
	return KeyService{
		keyRepository: repository,
		keyGenerator:  genRsaKey}, nil
}

var KeyAlreadyExists = errors.New("Key already exists")
