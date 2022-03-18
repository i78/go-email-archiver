package repository

import (
	"email-archiver-cli/mocks"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestMailRepository(t *testing.T) {
	mockController := gomock.NewController(t)

	t.Run("when new mail persist action requested", func(t *testing.T) {

		fakeLog := &log.Entry{}

		t.Run("should try to save to the expected filename", func(t *testing.T) {
			mockFilesystem := mocks.NewMockFilesystem(mockController)
			subject := MailRepository{mockFilesystem, fakeLog}

			mockFilesystem.EXPECT().Save("mailArchive/maildir/fake@fake.com/2000/01/01/01/fake", gomock.Any())

			mockContent := []byte{1, 2, 4, 8, 16}
			filenameHint := "fake@fake.com/2000/01/01/01/fake"

			subject.Persist(mockContent, filenameHint)
		})

	})
}
