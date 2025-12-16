package util

import (
	"log"
	"os"
)

type Logger struct {
	l *log.Logger
}

func NewLogger() *Logger {
	return &Logger{l: log.New(os.Stderr, "", log.LstdFlags)}
}

func (lg *Logger) Infof(format string, args ...any)  { lg.l.Printf("INFO "+format, args...) }
func (lg *Logger) Warnf(format string, args ...any)  { lg.l.Printf("WARN "+format, args...) }
func (lg *Logger) Errorf(format string, args ...any) { lg.l.Printf("ERROR "+format, args...) }
