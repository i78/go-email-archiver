package util

import log "github.com/sirupsen/logrus"

func PrepareLogging(fields log.Fields) *log.Entry {
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})
	log.SetLevel(log.DebugLevel)
	return log.WithFields(fields)
}
