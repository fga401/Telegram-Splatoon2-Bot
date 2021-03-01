package log

import (
	"go.uber.org/zap"
)

var logger *zap.Logger

func init() {
	logger, _ = zap.NewProduction()
}

// InitLogger initializes the logger with level.
// level should be 'debug', 'info', 'warn' or 'error'.
func InitLogger(level string) {
	config := zap.NewProductionConfig()
	switch level {
	case "debug":
		config.Level.SetLevel(zap.DebugLevel)
	case "info":
		config.Level.SetLevel(zap.InfoLevel)
	case "warn":
		config.Level.SetLevel(zap.WarnLevel)
	case "error":
		config.Level.SetLevel(zap.ErrorLevel)
	default:
		config.Level.SetLevel(zap.InfoLevel)
	}
	l, err := config.Build()
	if err != nil {
		Panic("can't initialize zap log", zap.Error(err))
	}
	logger = l
}

// Debug logs a message at debug level.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Info logs a message at info level.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Warn logs a message at warn level.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Error logs a message at error level.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Fatal logs a message at fatal level. The logger then calls os.Exit(1).
func Fatal(msg string, fields ...zap.Field) {
	logger.Fatal(msg, fields...)
}

// Panic logs a message at panic level. The logger then panics.
func Panic(msg string, fields ...zap.Field) {
	logger.Panic(msg, fields...)
}

// Sync flushes any buffered log entries.
func Sync() error {
	return logger.Sync()
}
