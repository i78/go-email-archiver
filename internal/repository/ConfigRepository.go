package repository

import (
	"email-archiver-cli/internal/filesystem"
	"email-archiver-cli/proto/emailarchver/config"
	"google.golang.org/protobuf/proto"
)

type FileConfigRepository struct {
	parentRepository Repository
}

func NewFileConfigRepository(parentRepository Repository) *FileConfigRepository {
	return &FileConfigRepository{parentRepository: parentRepository}
}

func (f *FileConfigRepository) Contains(configName string) bool {
	folderName := createFilename(configSubfolderPrefix, configName)
	return filesystem.Exists(folderName)
}

func (f *FileConfigRepository) Persist(mailConfig config.FetchMailConfig) error {
	fileName := createFilename(configSubfolderPrefix, mailConfig.Name)

	if wire, err := proto.Marshal(&mailConfig); err != nil {
		return err
	} else {
		err = f.parentRepository.filesystem.Save(fileName, wire)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FileConfigRepository) Get(mailConfigName string) (config config.FetchMailConfig, err error) {
	fileName := createFilename(configSubfolderPrefix, mailConfigName)

	if wire, err := f.parentRepository.filesystem.Read(fileName); err == nil {
		err = proto.Unmarshal(wire, &config)
	}

	return
}
