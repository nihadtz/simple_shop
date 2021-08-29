package services

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

var (
	Logger *logrus.Logger
)

func NewLogger() {
	fmt.Println("Initializing Logger")

	var log = logrus.New()

	file, err := os.OpenFile("logs/logs.log", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0777)

	if err != nil {
		fmt.Println("Error getting log file", err.Error())
	}

	log.Out = file
	Logger = log
}

func LogError(details string, err error) {
	_, filePath, line, _ := runtime.Caller(1)

	_, file := filepath.Split(filePath)

	Logger.WithField("file", file).WithField("line", line).Errorln(details, err.Error())
}

// LogInfo ->
func LogInfo(details string) {
	fmt.Println(details)

	_, filePath, line, _ := runtime.Caller(1)

	_, file := filepath.Split(filePath)

	Logger.WithField("file", file).WithField("line", line).Infoln(details)
}
