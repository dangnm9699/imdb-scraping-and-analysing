
package logger

import (
	"log"
	"os"
)

const logDir = "logs"

var (
	Info  *log.Logger
	Error *log.Logger
	Debug *log.Logger
)

func Set(filename string) {
	path := logDir + string(os.PathSeparator) + filename
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalln("failed to open logs.txt:", err)
	}
	Info = log.New(file, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(file, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(file, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
	Info.Println("======== START LOGGING ========")
}