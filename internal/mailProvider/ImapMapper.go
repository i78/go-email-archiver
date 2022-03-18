package mailProvider

import (
	"email-archiver-cli/proto/emailarchver/email"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/mail"
	"time"
)

func mapEMail(message *mail.Message, bodyReader func() ([]byte, error)) email.EMail {
	header := message.Header
	messageDate, _ := time.Parse(time.RFC1123Z, header.Get("Date"))

	headers := make(map[string]string)

	for k, v := range header {
		headers[k] = v[0]
	}

	body, _ := bodyReader()

	return email.EMail{
		MessageId:         header.Get("Message-Id"),
		Date:              timestamppb.New(messageDate),
		PreviousMessageId: "",
		Header:            &email.Header{Fields: headers},
		Parts: []*email.Part{
			{Raw: body},
		},
	}
}
