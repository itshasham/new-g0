package logger

import (
	"context"
	"log"
	"os"
	"runtime"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	defaultLogger *zap.Logger
	once          sync.Once
)

func init() {
	defaultLogger = NewLogger()
}

func NewLogger() *zap.Logger {
	return getLogger(zap.NewProductionConfig())
}

func NewTestLogger() *zap.Logger {
	return getLogger(zap.NewDevelopmentConfig())
}

func getLogger(cfg zap.Config) *zap.Logger {
	cfg.EncoderConfig.TimeKey = "timestamp"
	cfg.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)
	zapLogger, err := cfg.Build()
	if err != nil {
		log.Fatal(err)
	}
	return zapLogger
}

// logWithLevel logs a message with structured fields.
func logWithLevel(level zapcore.Level, ctx context.Context, message string, fields Fields) {
	// Always enrich logs with the correlation ID.
	allFields := make(Fields, len(fields)+1)
	for k, v := range fields {
		allFields[k] = v
	}
	allFields[FieldCorrelationID] = extractCorrelationID(ctx)

	// Convert the map to a slice of zap.Field
	zapFields := make([]zap.Field, 0, len(allFields))
	for k, v := range allFields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	// Create stack trace
	// Get caller information
	if pc, file, line, ok := runtime.Caller(2); ok {
		fn := runtime.FuncForPC(pc)
		if fn != nil {
			zapFields = append(zapFields, zap.String("caller", fn.Name()))
			zapFields = append(zapFields, zap.String("file", file))
			zapFields = append(zapFields, zap.Int("line", line))
		}
	}

	// Log with the structured fields.
	switch level {
	case zapcore.DebugLevel:
		defaultLogger.Debug(message, zapFields...)
	case zapcore.InfoLevel:
		defaultLogger.Info(message, zapFields...)
	case zapcore.WarnLevel:
		defaultLogger.Warn(message, zapFields...)
	case zapcore.ErrorLevel:
		defaultLogger.Error(message, zapFields...)
	case zapcore.FatalLevel:
		defaultLogger.Fatal(message, zapFields...)
	}
}

// Error logs an error message with structured fields
func Error(ctx context.Context, message string, fields Fields) {
	logWithLevel(zapcore.ErrorLevel, ctx, message, fields)
}

// Info logs an informational message with structured fields
func Info(ctx context.Context, message string, fields Fields) {
	logWithLevel(zapcore.InfoLevel, ctx, message, fields)
}

// Warn logs a warning message with structured fields
func Warn(ctx context.Context, message string, fields Fields) {
	logWithLevel(zapcore.WarnLevel, ctx, message, fields)
}

// Debug logs a debug message with structured fields
func Debug(ctx context.Context, message string, fields Fields) {
	logWithLevel(zapcore.DebugLevel, ctx, message, fields)
}

// Fatal logs an error message and exits cleanly without stack traces
func Fatal(ctx context.Context, message string, fields Fields) {
	logWithLevel(zapcore.ErrorLevel, ctx, message, fields)
	os.Exit(1)
}
