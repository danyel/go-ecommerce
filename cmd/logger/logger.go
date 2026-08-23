package logger

import (
	InternalLog "log"
	OS "os"
)

var Log Logger = &loggerImpl{}

type Logger interface {
	Debug(format string, args ...any)
	Info(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type loggerImpl struct {
}

func (logger *loggerImpl) Debug(format string, args ...any) {
	if OS.Getenv("DEBUG_ENABLED") == "true" {
		InternalLog.Printf("[DEBUG] "+format, args...)
	}
}

func (logger *loggerImpl) Info(format string, args ...any) {
	InternalLog.Printf("[INFO] "+format, args...)
}

func (logger *loggerImpl) Fatal(args ...any) {
	logger.Fatalf("[FATAL] %v", args...)
}

// Fatalf is equivalent to [log.Fatalf] followed by a call to [os.Exit](1).
func (logger *loggerImpl) Fatalf(format string, args ...any) {
	InternalLog.Fatalf("[FATAL] %s"+format, args...)
}
