package logger

import (
	"context"
	Fmt "fmt"
	SLog "log/slog"
	OS "os"
)

var (
	Log  Logger = &loggerImpl{}
	sLog        = SLog.New(SLog.NewTextHandler(OS.Stdout, &SLog.HandlerOptions{Level: GetLevel()}))
)

type ctxKey string

const CorrelationIDKey ctxKey = "correlation_id"

func GetLevel() SLog.Leveler {
	if OS.Getenv("LOG_LEVEL") != "DEBUG" {
		return SLog.LevelDebug
	}
	return SLog.LevelInfo
}

type Logger interface {
	Debug(format string, args ...any)
	DebugCtx(ctx context.Context, format string, args ...any)
	Info(format string, args ...any)
	Fatal(args ...any)
	Fatalf(format string, args ...any)
}

type loggerImpl struct {
}

func (logger *loggerImpl) Debug(format string, args ...any) {
	logger.getSlogWithCorrelation(context.Background()).Debug(Fmt.Sprintf(format, args...))
}
func (logger *loggerImpl) DebugCtx(ctx context.Context, format string, args ...any) {
	logger.getSlogWithCorrelation(ctx).Debug(Fmt.Sprintf(format, args...))
}

func (logger *loggerImpl) Info(format string, args ...any) {
	logger.getSlogWithCorrelation(context.Background()).Info(Fmt.Sprintf(format, args...))
}

func (logger *loggerImpl) Fatal(args ...any) {
	logger.Fatalf("%v", args...)
}

// Fatalf is equivalent to [log.Fatalf] followed by a call to [os.Exit](1).
func (logger *loggerImpl) Fatalf(format string, args ...any) {
	logger.getSlogWithCorrelation(context.Background()).Error(Fmt.Sprintf(format, args...))
	OS.Exit(0)
}

func (logger *loggerImpl) getSlogWithCorrelation(ctx context.Context) *SLog.Logger {
	if ctx != nil {
		if id, ok := ctx.Value(CorrelationIDKey).(string); ok && id != "" {
			return sLog.With("correlation_id", id)
		}
	}
	return sLog
}
