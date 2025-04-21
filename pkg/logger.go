package pkg

import (
	"log"
	"os"
)

type Logger struct {
	info *log.Logger
}

func (l *Logger) Info(s string) {
	l.info.Println(s)
}

func (l *Logger) Infof(format string, args ...any) {
	l.info.Printf(format, args...)
}

func GetLogger() *Logger {
	return &Logger{
		info: log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}
