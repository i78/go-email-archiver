package packager

import (
	"email-archiver-cli/internal/cryptography"
	"email-archiver-cli/models"
	"email-archiver-cli/proto/emailarchver/archive"
	"email-archiver-cli/proto/emailarchver/email"
	"google.golang.org/protobuf/proto"
	"time"
)

type Packager struct {
	Key models.Key
}

func (p *Packager) sign(mail email.EMail, key models.Key) (archive.ArchiveEnvelope, error) {
	mailWire, _ := proto.Marshal(&mail)
	mailWireEncrypted, _ := cryptography.Encrypt(&key, mailWire) // todo: this should be locallly injected.

	hash, _ := cryptography.Sign(&key, mailWire)

	envelope := archive.ArchiveEnvelope{
		KeyIdHint: key.Revision.String(),
		Content:   mailWireEncrypted,
		Log: []*archive.SignedLogEntry{
			{
				Time:            time.Now().Unix(),
				Principal:       "",
				Event:           archive.EventType_CREATED,
				ContentChecksum: nil,
				Signature:       hash,
			},
		},
	}
	return envelope, nil
}

func (p *Packager) Package(email email.EMail) (archive.ArchiveEnvelope, error) {
	return p.sign(email, p.Key)
}
