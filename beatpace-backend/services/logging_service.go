package services

import (
	"fmt"
	"log"
	"os"
)

// LogError logs an error message
func LogError(err error) {
	if err != nil {
		log.Printf("Error: %v", err)
	}
}

// LogInfo logs information messages (for example, to track the flow of requests)
func LogInfo(message string) {
	log.Printf("Info: %s", message)
}

// Setup logging to a file
func InitLogging() {
	logFile, err := os.OpenFile("app.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Error opening log file:", err)
		return
	}
	log.SetOutput(logFile)
}
