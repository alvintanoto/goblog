package log

import (
	"log"
	"os"
)

type Log struct {
	ErrorLog *log.Logger
	InfoLog  *log.Logger
}

func Get() *Log {
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	return &Log{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}
}
