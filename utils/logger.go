package utils

import (
    "log"
    "os"
)

var (
    InfoLogger  *log.Logger
    ErrorLogger *log.Logger
)

func InitLogger() {
    file, err := os.OpenFile("logs/app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
    if err != nil {
        log.Fatalln("Failed to open log file:", err)
    }
    InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
    ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
}
