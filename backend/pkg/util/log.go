package util

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger *zap.Logger
	once   sync.Once
)

// InitLogger initializes the zap logger with the given configuration.
func InitLogger(logToFile bool, filePath string) {
	once.Do(func() {
		var core zapcore.Core

		// Define encoder configuration
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "time",
			LevelKey:       "level",
			NameKey:        "logger",
			CallerKey:      "caller",
			MessageKey:     "msg",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.StringDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}

		// Create console encoder
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

		// Create log output targets
		var writeSyncers []zapcore.WriteSyncer
		writeSyncers = append(writeSyncers, zapcore.AddSync(os.Stdout))

		if logToFile {
			file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				panic(err)
			}
			writeSyncers = append(writeSyncers, zapcore.AddSync(file))
		}

		// Combine write syncers
		multiWriteSyncer := zapcore.NewMultiWriteSyncer(writeSyncers...)

		// Create core
		core = zapcore.NewCore(consoleEncoder, multiWriteSyncer, zapcore.DebugLevel)

		// Create logger
		logger = zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))
	})
}

// Info logs an info message.
func Info(msg string, fields ...zap.Field) {
	logger.Info(msg, fields...)
}

// Error logs an error message.
func Error(msg string, fields ...zap.Field) {
	logger.Error(msg, fields...)
}

// Debug logs a debug message.
func Debug(msg string, fields ...zap.Field) {
	logger.Debug(msg, fields...)
}

// Warn logs a warning message.
func Warn(msg string, fields ...zap.Field) {
	logger.Warn(msg, fields...)
}

// Sync flushes any buffered log entries.
func Sync() {
	_ = logger.Sync()
}