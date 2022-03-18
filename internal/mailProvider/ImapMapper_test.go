package mailProvider

import (
	"email-archiver-cli/proto/emailarchver/email"
	"github.com/stretchr/testify/assert"
	"net/mail"
	"testing"
)

func TestIMAPMapper(t *testing.T) {
	t.Run("eMail Mapper", func(t *testing.T) {
		fakeMessageId := "cafebabe-chronologie-4@dreese.de"
		fakeMailHeader := mail.Header{
			"Date":              []string{"Sun, 13 Feb 2022 18:04:16 +0100"},
			"Message-Id":        []string{fakeMessageId},
			"fake_header_0_key": []string{"fake_header_0_value"},
		}

		fakeMail := mail.Message{
			Header: fakeMailHeader,
			Body:   nil,
		}

		t.Run("should contain correctly mapped date", func(t *testing.T) {
			result := mapEMail(&fakeMail, fakeBodyReader)
			assert.NotNil(t, result.Date)
		})

		t.Run("should contain all original headers", func(t *testing.T) {
			result := mapEMail(&fakeMail, fakeBodyReader)

			expectedHeaders := &email.Header{Fields: map[string]string{
				"Date":              "Sun, 13 Feb 2022 18:04:16 +0100",
				"Message-Id":        fakeMessageId,
				"fake_header_0_key": "fake_header_0_value",
			}}

			assert.Equal(t, expectedHeaders, result.Header)
		})

		t.Run("should contain all original message ID", func(t *testing.T) {
			result := mapEMail(&fakeMail, fakeBodyReader)

			assert.Equal(t, fakeMessageId, result.MessageId)
		})
	})
}

func fakeBodyReader() (content []byte, err error) {
	return []byte{0, 1, 2, 3}, nil
}
