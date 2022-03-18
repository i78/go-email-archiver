package repository

import (
	"crypto/x509"
	"email-archiver-cli/internal/filesystem"
	"email-archiver-cli/internal/util"
	"email-archiver-cli/models"
	"github.com/google/uuid"
	log "github.com/sirupsen/logrus"
	"io/fs"
)

type FileKeyRepository struct {
	filesystem Filesystem
	log        *log.Entry
}

func NewFileKeyRepository(parentRepository Repository) *FileKeyRepository {
	return &FileKeyRepository{
		filesystem: parentRepository.filesystem,
		log:        util.PrepareLogging(log.Fields{"source": "FileKeyRepository"})}
}

func (f *FileKeyRepository) Contains(keyName string) bool {
	folderName := createFilename(keySubfolderPrefix, keyName)
	return filesystem.Exists(folderName)
}

func (f *FileKeyRepository) Persist(keyName string, key *models.Key) error {
	keyRootFolder := []string{keySubfolderPrefix, keyName}
	publicKeyFolder := append(keyRootFolder, "public")
	privateKeyFolder := append(keyRootFolder, "private")

	for _, folder := range [][]string{keyRootFolder, publicKeyFolder, privateKeyFolder} {
		filename := createFilename(folder...)
		if err := f.filesystem.MkDir(filename); err != nil {
			return err
		}
	}

	filenameWithRevision := func(path []string) string {
		return createFilename(append(path, key.Revision.String())...)
	}

	privateKeyFileName, publicKeyFileName :=
		filenameWithRevision(privateKeyFolder),
		filenameWithRevision(publicKeyFolder)

	pub := x509.MarshalPKCS1PublicKey(key.PublicKey)
	errPub := f.filesystem.Save(publicKeyFileName, pub)
	if errPub != nil {
		return errPub
	}

	prv := x509.MarshalPKCS1PrivateKey(key.PrivateKey)

	errPk := f.filesystem.Save(privateKeyFileName, prv)
	if errPk != nil {
		return errPk
	}

	return nil
}

func (f *FileKeyRepository) Get(keyName string) (key *models.Key, err error) {
	folderName := createFilename(keySubfolderPrefix, keyName, "private")

	revisionUuid, _ := f.getLatestKeyRevision(folderName)

	f.log.WithFields(log.Fields{
		"keyName":  keyName,
		"revision": revisionUuid.String(),
	}).Info()

	privateKeyFileName := createFilename(keySubfolderPrefix, keyName, "private", revisionUuid.String())

	privateKeyBytes, fileReadError := f.filesystem.Read(privateKeyFileName)
	if fileReadError != nil {
		log.WithFields(log.Fields{
			"keyName":       keyName,
			"revision":      revisionUuid.String(),
			"fileReadError": fileReadError,
		}).Fatal("unable to read key from filesystem!")
		return nil, fileReadError
	}

	privateKey, keyParseError := x509.ParsePKCS1PrivateKey(privateKeyBytes)
	if keyParseError != nil {
		f.log.WithFields(log.Fields{
			"keyName":    keyName,
			"revision":   revisionUuid.String(),
			"parseError": keyParseError,
		}).Fatal("unable to parse key!")
		return nil, keyParseError
	}

	key = &models.Key{
		Revision:   uuid.UUID{},
		PublicKey:  &privateKey.PublicKey,
		PrivateKey: privateKey,
	}

	return
}

func (f *FileKeyRepository) getLatestKeyRevision(folderName string) (uuid.UUID, error) {
	files, _ := f.filesystem.Ls(folderName)

	latestKey := util.Max(&files, func(o interface{}, p interface{}) bool {
		return o.(fs.FileInfo).ModTime().Unix() < p.(fs.FileInfo).ModTime().Unix()
	})

	revisionUuid, _ := uuid.Parse(latestKey.Name())
	return revisionUuid, nil
}
