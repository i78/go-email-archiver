package services

import "email-archiver-cli/proto/emailarchver/config"

type ConfigRepository interface {
	Contains(configName string) bool
	Persist(configName string, config config.FetchMailConfig) error
}
