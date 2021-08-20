package logger

import (
	"log"
	"os"
)

var (
	// Info is a logger that record information
	Info  *log.Logger
	// Error is a logger that record errors
	Error *log.Logger
	// Debug is a logger that record stuff that might be helpful to debug
	Debug *log.Logger
)

// Set initializes loggers
func Set() {
	Info = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
	Debug = log.New(os.Stdout, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
}
