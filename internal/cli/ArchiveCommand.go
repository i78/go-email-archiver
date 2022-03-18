package cli

import (
	context2 "context"
	"email-archiver-cli/internal/archiver"
	"email-archiver-cli/internal/mailProvider"
	"email-archiver-cli/internal/packager"
	"email-archiver-cli/internal/repository"
	"email-archiver-cli/internal/util"
	"email-archiver-cli/models"
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

type ArchiveCommand struct {
	ConfigName string
	KeyName    string
	log        *log.Entry
}

func (sv *ArchiveCommand) Run(_ *kong.Context) (err error) {
	sv.log = util.PrepareLogging(log.Fields{"command": "ARCHIVE", "configName": sv.ConfigName})
	sv.log.Info("Starting email archive process")

	configRepo, _ := sv.BuildConfigRepo()
	if config, getConfigError := configRepo.Get(sv.ConfigName); getConfigError != nil {
		sv.log.Fatal("Config does not exist", getConfigError)
		err = getConfigError
	} else {
		sv.log.Info("loaded config")
		key, keyError := sv.resolveKey(sv.KeyName)
		if keyError != nil {
			sv.log.Fatal("Unable to load key", keyError)
		}

		pk := packager.Packager{
			Key: *key,
		}

		newMailSink, _ := sv.BuildMailRepo()

		newMailHandler := archiver.NewArchiver(&pk, newMailSink, config.Name)

		context, _ := context2.WithCancel(context2.Background())
		mailProvider.NewIMAPProvider(context, config, newMailHandler.NewMail)

		sv.log.Info("Started email archive process")
	}

	return
}

func (sv *ArchiveCommand) BuildConfigRepo() (configRepo *repository.FileConfigRepository, err error) {
	if repo, err := repository.Open(); err != nil {
		log.Fatal("Unable to open repository.")
	} else {
		configRepo = repository.NewFileConfigRepository(*repo)
	}
	return
}

func (sv *ArchiveCommand) BuildMailRepo() (configRepo *repository.MailRepository, err error) {
	if repo, err := repository.Open(); err != nil {
		log.Fatal("Unable to open repository.")
	} else {
		configRepo = repository.NewMailRepository(*repo)
	}
	return
}

func (sv *ArchiveCommand) resolveKey(keyName string) (key *models.Key, err error) {
	err = sv.withRootRepository(func(mailRepository *repository.Repository) error {
		keyRepo := repository.NewFileKeyRepository(*mailRepository)

		key, err = keyRepo.Get(keyName)

		return nil
	})
	return
}

func (sv *ArchiveCommand) withRootRepository(fn func(*repository.Repository) error) error {
	repo, err := repository.Open()
	if err != nil {
		log.Fatal("Unable to open repository.")
	} else {
		return fn(repo)
	}
	return nil
}
