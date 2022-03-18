package archiver

type Sink interface {
	Persist(content []byte, filenameHint string) error
}
