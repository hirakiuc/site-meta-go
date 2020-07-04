package logger

import (
	"io/ioutil"
	"log"
	"os"
)

func isEnabledLogger() bool {
	return len(os.Getenv("LOG_ENABLE")) > 0
}

// GetLogger return private logger instance.
func GetLogger() *log.Logger {
	instance := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)

	if !isEnabledLogger() {
		instance.SetOutput(ioutil.Discard)
	}

	return instance
}
