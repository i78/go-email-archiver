package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"email-archiver-cli/models"
)

func Encrypt(key *models.Key, content []byte) (encrypted []byte, err error) {
	return rsa.EncryptPKCS1v15(rand.Reader, &key.PrivateKey.PublicKey, content)
}
