package logger

import (
	"fmt"
	"log"
	"os"
)

const logsDirPath = "/var/log/patcher-container-runtime.log"

func getLogFilePath(appName string) string {
	return fmt.Sprintf("%v-%v", logsDirPath, appName)
}

func New(appName string) *log.Logger {
	var logger *log.Logger
	logFile, err := os.OpenFile(getLogFilePath(appName), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err == nil {
		logger = log.New(logFile, fmt.Sprintf("[%v]", appName), log.LstdFlags)
	} else {
		logger = log.Default()
	}

	logger.Printf(fmt.Sprintf("Logger %v initialized", appName))
	return logger
}
