package cryptography

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha512"
	"email-archiver-cli/models"
)

func Sign(key *models.Key, content []byte) (signature []byte, err error) {
	hasher := sha512.New()
	hasher.Write(content)
	hash := hasher.Sum(nil)

	signature, err = rsa.SignPKCS1v15(rand.Reader, key.PrivateKey, crypto.SHA512, hash)
	return
}

func IsValid(key *models.Key, content []byte, signature []byte) (result bool) {
	hasher := sha512.New()
	hasher.Write(content)
	hash := hasher.Sum(nil)

	result = rsa.VerifyPKCS1v15(key.PublicKey, crypto.SHA512, hash, signature) == nil
	return
}
