package sitemeta

import (
	"io/ioutil"
	"log"
	"os"
)

var logger *log.Logger

func isEnabledLogger() bool {
	return len(os.Getenv("LOG_ENABLE")) > 0
}

// SetLogger set user's logger to private logger.
func SetLogger(newLogger *log.Logger) {
	logger = newLogger
}

func initLogger() {
	logger = log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	if isEnabledLogger() == false {
		logger.SetOutput(ioutil.Discard)
	}
}
