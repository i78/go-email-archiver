package archiver

//go:generate  mockgen -source Packager.go -destination ../../mocks/Packager.go -package=mocks
//go:generate  mockgen -source Sink.go -destination ../../mocks/Sink.go -package=mocks

import (
	"email-archiver-cli/mocks"
	"email-archiver-cli/proto/emailarchver/archive"
	"email-archiver-cli/proto/emailarchver/email"
	"github.com/golang/mock/gomock"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestArchiver(t *testing.T) {
	mockController := gomock.NewController(t)
	mockPackager := mocks.NewMockPackager(mockController)
	mockSink := mocks.NewMockSink(mockController)

	const fakeConfigName = "fakeConfigName"

	subject := Archiver{
		configName: fakeConfigName,
		packager:   mockPackager,
		sink:       mockSink,
		log:        log.WithFields(log.Fields{}),
	}

	t.Run("when new mail provided", func(t *testing.T) {

		fakeMail := email.EMail{
			MessageId:         "<18e4e683a8cf769491574ce16a4ea13b8fdf44b7.camel@dreese.de>",
			PreviousMessageId: "<bb3dc1b58804202c8a700deecb826c880a8af60b.camel@dreese.de>",
			Date:              timestamppb.New(time.Date(2000, 01, 01, 00, 00, 00, 0, time.UTC)),
			Header:            nil,
			Parts:             nil,
		}

		const expectedFileName = "fakeConfigName/2000/01/01/18e4e683a8cf769491574ce16a4ea13b8fdf44b7.camel@dreese.de"

		t.Run("should encrypt and sign mail", func(t *testing.T) {
			mockPackager.EXPECT().Package(fakeMail).Return(archive.ArchiveEnvelope{}, nil)
			mockSink.EXPECT().Persist(gomock.Any(), expectedFileName).Return(nil) // no default mocks? ORLY?
			subject.NewMail(fakeMail)
		})

		t.Run("should generate expected filename", func(t *testing.T) {
			result := generateFileName(fakeConfigName, &fakeMail)
			assert.Equal(t, expectedFileName, result)
		})

		t.Run("should try to persist email to expected directory", func(t *testing.T) {
			mockPackager.EXPECT().Package(fakeMail).Return(archive.ArchiveEnvelope{}, nil)
			mockSink.EXPECT().Persist(gomock.Any(), expectedFileName).Return(nil)

			subject.NewMail(fakeMail)
		})

	})
}
