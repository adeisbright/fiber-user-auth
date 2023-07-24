package common

import (
	"log"
	"os"
)

//Sends a log message to the logs directory using the message
// and the level provided

func CustomLog(message string, level string) error {

	filename := "./logs/log.txt"
	file, err := os.OpenFile(filename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	logger := log.New(file, level+": ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Println(message)
	return nil

}
