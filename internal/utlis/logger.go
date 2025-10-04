package utils

import (
	"io"
	"log"
	"os"
)

var (
	Info  *log.Logger
	Error *log.Logger
)

func InitLogger(logFilePath string) {
	file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file %s: %v", logFilePath, err)
	}

	multiWriter := io.MultiWriter(file, os.Stdout)

	Info = log.New(multiWriter, "📥 [INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(multiWriter, "❌ [ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
}
