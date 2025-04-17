package utils

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Logger interface {
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Critical(format string, v ...interface{})
}

type logger struct {
	logLevel int
	logger   *log.Logger
}

const (
	LevelCritical = iota
	LevelError
	LevelWarning
	LevelInfo
)

func NewLogger(level string) Logger {
	var logLevel int
	switch strings.ToLower(level) {
	case "critical":
		logLevel = LevelCritical
	case "error":
		logLevel = LevelError
	case "warning":
		logLevel = LevelWarning
	case "info":
		logLevel = LevelInfo
	default:
		logLevel = LevelInfo
	}

	return &logger{
		logLevel: logLevel,
		logger:   log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	}
}

func (l *logger) Info(format string, v ...interface{}) {
	if l.logLevel >= LevelInfo {
		l.logger.Output(2, fmt.Sprintf("INFO: "+format, v...))
	}
}

func (l *logger) Warning(format string, v ...interface{}) {
	if l.logLevel >= LevelWarning {
		l.logger.Output(2, fmt.Sprintf("WARNING: "+format, v...))
	}
}

func (l *logger) Error(format string, v ...interface{}) {
	if l.logLevel >= LevelError {
		l.logger.Output(2, fmt.Sprintf("ERROR: "+format, v...))
	}
}

func (l *logger) Critical(format string, v ...interface{}) {
	if l.logLevel >= LevelCritical {
		l.logger.Output(2, fmt.Sprintf("CRITICAL: "+format, v...))
		os.Exit(1)
	}
}
