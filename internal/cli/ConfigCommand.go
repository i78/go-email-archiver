package cli

import (
	"email-archiver-cli/internal/repository"
	"email-archiver-cli/proto/emailarchver/config"
	"github.com/alecthomas/kong"
	log "github.com/sirupsen/logrus"
)

type ConfigCommand struct {
	Create Create `cmd:"" help:"Create a new configuration"`
	Delete struct {
		ConfigName string
	} `cmd:"" help:"Delete a configuration"`
	List struct {
	} `cmd:"" help:"List existing configurations"`
}

type Create struct {
	Name          string
	Username      string
	Password      string
	Server        string
	CheckInterval uint32
}

func (sv *Create) Run(ctx *kong.Context) error {
	log.WithFields(log.Fields{"k": ctx.Selected()}).Info("Saving new Config")

	if sv.Username != "" {
		repo, repoErr := repository.Open()

		if repoErr != nil {
			log.Fatal("Unable to open repository.")
			return repoErr
		}

		configRepo := repository.NewFileConfigRepository(*repo)

		newConfig := sv.mapToFetchMailConfig()

		configRepo.Persist(newConfig)
	}

	return nil
}

func (sv *Create) mapToFetchMailConfig() config.FetchMailConfig {
	return config.FetchMailConfig{
		Name:          sv.Name,
		Username:      sv.Username,
		Password:      sv.Password,
		CheckInterval: sv.CheckInterval,
		Server:        sv.Server,
	}
}
