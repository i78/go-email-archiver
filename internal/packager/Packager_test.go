package packager

import (
	"crypto/rand"
	"crypto/rsa"
	"email-archiver-cli/models"
	"email-archiver-cli/proto/emailarchver/archive"
	"email-archiver-cli/proto/emailarchver/email"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestPackager(t *testing.T) {

	t.Run("when new mail provided", func(t *testing.T) {
		subject := Packager{}
		fakeKey, _ := genFakeKey()

		fakeMail := email.EMail{
			MessageId:         "FAKE_MSG_ID",
			PreviousMessageId: "FAKE_PREVIOUS_MSG_ID",
			Header:            nil,
			Parts:             nil,
		}

		t.Run("sign function", func(t *testing.T) {

			signedLogEntry, _ := subject.sign(fakeMail, *fakeKey)

			t.Run("should map to log entry", func(t *testing.T) {
				assert.NotNil(t, signedLogEntry)
			})

			t.Run("should contain content neq nil", func(t *testing.T) {
				assert.NotNil(t, signedLogEntry.Content)
			})

			t.Run("should have hint of key used to encrypt the body", func(t *testing.T) {
				assert.Equal(t, signedLogEntry.KeyIdHint, fakeKey.Revision.String())
			})

			t.Run("should have a signature stating that the mail has been archived recently", func(t *testing.T) {
				assert.Equal(t, len(signedLogEntry.Log), 1)
				assert.Equal(t, signedLogEntry.Log[0].Event, archive.EventType_CREATED)
			})

			t.Run("should contain a checksum for the content", func(t *testing.T) {
				assert.NotNil(t, signedLogEntry.Log[0].Signature)
			})

		})
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
