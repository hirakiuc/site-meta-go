package logger

import (
	"io/ioutil"
	"log"
	"os"
)

var logger = newLogger()

func isEnabledLogger() bool {
	return len(os.Getenv("LOG_ENABLE")) > 0
}

func newLogger() *log.Logger {
	instance := log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	if isEnabledLogger() == false {
		instance.SetOutput(ioutil.Discard)
	}
	return instance
}

// GetLogger return private logger instance.
func GetLogger() *log.Logger {
	return logger
}
