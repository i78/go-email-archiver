package cryptography

import (
	"crypto/rand"
	"crypto/rsa"
	"email-archiver-cli/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSigner(t *testing.T) {

	t.Run("should return signature for given content", func(t *testing.T) {
		fakeKey, _ := genFakeKey()
		mockContent := []byte{1, 2, 4, 8, 16}

		signature, err := Sign(fakeKey, mockContent)

		assert.Nil(t, err)
		assert.NotEmpty(t, signature)
	})

	t.Run("should return positive validation result when validated right away", func(t *testing.T) {
		fakeKey, _ := genFakeKey()
		mockContent := []byte{1, 2, 4, 8, 16}

		signature, _ := Sign(fakeKey, mockContent)
		isValid := IsValid(fakeKey, mockContent, signature)

		assert.True(t, isValid)
	})

	t.Run("should return negative validation result when content has been tampered with", func(t *testing.T) {
		fakeKey, _ := genFakeKey()
		mockContent := []byte{1, 2, 4, 8, 16}
		tamperedMockContent := []byte{1, 2, 4, 8, 8}

		signature, _ := Sign(fakeKey, mockContent)
		isValid := IsValid(fakeKey, tamperedMockContent, signature)

		assert.False(t, isValid)
	})
}

func genFakeKey() (key *models.Key, err error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, err
	}
	publicKey := &privateKey.PublicKey

	revision, errUuid := uuid.NewUUID()
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
