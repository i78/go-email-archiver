package repository

//go:generate  mockgen -source Filesystem.go -destination ../../mocks/Filesystem.go -package=mocks

import (
	"crypto/rand"
	"crypto/rsa"
	"email-archiver-cli/mocks"
	"email-archiver-cli/models"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFileKeyRepository(t *testing.T) {
	mockController := gomock.NewController(t)

	t.Run("when new key persisted", func(t *testing.T) {

		fakeLog := &log.Entry{}
		fakeKey, _ := genFakeKey()
		fmt.Println("go")

		t.Run("should ensure that the rootfolder in parent repo exists", func(t *testing.T) {
			mockFilesystem := mocks.NewMockFilesystem(mockController)
			subject := FileKeyRepository{mockFilesystem, fakeLog}

			mockFilesystem.EXPECT().MkDir("mailArchive/keys/fake-key")
			mockFilesystem.EXPECT().MkDir("mailArchive/keys/fake-key/public")
			mockFilesystem.EXPECT().MkDir("mailArchive/keys/fake-key/private")
			mockFilesystem.EXPECT().Save(gomock.Any(), gomock.Any()).AnyTimes()

			subject.Persist("fake-key", fakeKey)
		})

		t.Run("should throw error when unable to create directories", func(t *testing.T) {
			mockFilesystem := mocks.NewMockFilesystem(mockController)
			subject := FileKeyRepository{mockFilesystem, fakeLog}

			mockFilesystem.EXPECT().MkDir(gomock.Any()).Return(fmt.Errorf("boom"))

			err := subject.Persist("fake-key", fakeKey)

			assert.Error(t, err)
		})

		t.Run("should try to save keys to storage", func(t *testing.T) {
			mockFilesystem := mocks.NewMockFilesystem(mockController)
			subject := FileKeyRepository{mockFilesystem, fakeLog}

			mockFilesystem.EXPECT().MkDir(gomock.Any()).AnyTimes()
			mockFilesystem.EXPECT().Save("mailArchive/keys/fake-key/public/"+fakeKey.Revision.String(), gomock.Any())
			mockFilesystem.EXPECT().Save("mailArchive/keys/fake-key/private/"+fakeKey.Revision.String(), gomock.Any())

			subject.Persist("fake-key", fakeKey)
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
