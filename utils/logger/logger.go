package logger

import (
	"log"
	"os"
)

type Logger struct {
	err  *log.Logger
	warn *log.Logger
}

func New() Logger {
	return Logger{
		err:  log.New(os.Stderr, "ERR:  ", log.LstdFlags),
		warn: log.New(os.Stdout, "WARN: ", log.LstdFlags),
	}
}

func (logger *Logger) Warningf(fromat string, args ...any) {
	logger.warn.Printf(fromat, args...)
	logger.warn.Println()
}

func (logger *Logger) Errorf(fromat string, args ...any) {
	logger.err.Printf(fromat, args...)
	logger.warn.Println()
}

func (logger *Logger) Warning(err error) {
	logger.warn.Println(err.Error())
}

func (logger *Logger) Error(err error) {
	logger.err.Println(err.Error())
}