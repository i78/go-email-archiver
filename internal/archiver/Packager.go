package archiver

import (
	"email-archiver-cli/proto/emailarchver/archive"
	"email-archiver-cli/proto/emailarchver/email"
)

type Packager interface {
	Package(mail email.EMail) (archive.ArchiveEnvelope, error)
}
