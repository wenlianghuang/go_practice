package main

import (
	sublogrus "ToReadMeLog/Sublogrus"
	"os"
)

func main() {
	LogFilePath, errLogFilePath := os.Getwd()
	if errLogFilePath != nil {
		panic(errLogFilePath)
	}

	log := sublogrus.Sublogrusfunc(LogFilePath)
	log.Info("Record to the README.md")
}
