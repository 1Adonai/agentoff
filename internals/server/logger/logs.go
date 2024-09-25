package logger

import (
	"io"
	"log"
	"os"
)

var logFile *os.File

func InitLogger() {
	var err error
	logFile, err = os.OpenFile("./internals/server/logger/logs.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	multiWriter := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multiWriter)
}

func CloseLogger() {
	if logFile != nil {
		logFile.Close()
	}
}
