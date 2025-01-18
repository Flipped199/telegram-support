package log

import (
	"log/slog"
	"os"
)

type Logger interface {
	Error(msg string, err error, keysAndValues ...any)
	Info(msg string, keysAndValues ...any)
	Debug(msg string, keysAndValues ...any)
	Warn(msg string, keysAndValues ...any)
}

// SlogLogger 是基于 slog 的实现
type SlogLogger struct {
	logger *slog.Logger
}

var defaultLogger = NewSlogLogger()

var logLevel = new(slog.LevelVar)

func NewSlogLogger() *SlogLogger {
	return &SlogLogger{
		logger: slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: logLevel,
		})),
	}
}

func (l *SlogLogger) Error(msg string, err error, keysAndValues ...any) {
	if err != nil {
		keysAndValues = append(keysAndValues, slog.String("error", err.Error()))
	}
	l.logger.Error(msg, keysAndValues...)
}

func (l *SlogLogger) Info(msg string, keysAndValues ...any) {
	l.logger.Info(msg, keysAndValues...)
}

func (l *SlogLogger) Warn(msg string, keysAndValues ...any) {
	l.logger.Warn(msg, keysAndValues...)
}

func (l *SlogLogger) Debug(msg string, keysAndValues ...any) {
	l.logger.Debug(msg, keysAndValues...)
}

func Info(msg string, keysAndValues ...any) {
	defaultLogger.Info(msg, keysAndValues...)
}

func Warn(msg string, keysAndValues ...any) {
	defaultLogger.Warn(msg, keysAndValues...)
}

func Error(msg string, err error, keysAndValues ...any) {
	defaultLogger.Error(msg, err, keysAndValues...)
}

func Debug(msg string, keysAndValues ...any) {
	defaultLogger.Debug(msg, keysAndValues...)
}

func Fatal(msg string, err error, keysAndValues ...any) {
	defaultLogger.Error(msg, err, keysAndValues...)
	os.Exit(1)
}

func SetLevel(l slog.Level) {
	logLevel.Set(l)
}
