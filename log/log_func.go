package log

import (
	"context"

	"go.uber.org/zap"
)

func Debug(msg string, fields ...zap.Field) {
	_log.Debug(msg, fields...)
}

func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Debug(msg, fields...)
}

func Info(msg string, fields ...zap.Field) {
	_log.Info(msg, fields...)
}

func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Info(msg, fields...)
}

func Warn(msg string, fields ...zap.Field) {
	_log.Warn(msg, fields...)
}

func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Warn(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	_log.Error(msg, fields...)
}

func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Error(msg, fields...)
}

func Panic(msg string, fields ...zap.Field) {
	_log.Panic(msg, fields...)
}

func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Panic(msg, fields...)
}

func Fatal(msg string, fields ...zap.Field) {
	_log.Fatal(msg, fields...)
}

func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	GetLogger(ctx).Fatal(msg, fields...)
}

func Debugf(format string, args ...interface{}) {
	_log.Debugf(format, args...)
}

func DebugfCtx(ctx context.Context, format string, args ...interface{}) {
	GetLogger(ctx).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	_log.Infof(format, args...)
}

func InfofCtx(ctx context.Context, format string, args ...interface{}) {
	GetLogger(ctx).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	_log.Warnf(format, args...)
}

func WarnfCtx(ctx context.Context, format string, args ...interface{}) {
	GetLogger(ctx).Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	_log.Errorf(format, args...)
}

func ErrorfCtx(ctx context.Context, format string, args ...interface{}) {
	GetLogger(ctx).Errorf(format, args...)
}

func Panicf(format string, args ...interface{}) {
	_log.Panicf(format, args...)
}

func PanicfCtx(ctx context.Context, format string, args ...interface{}) {
	GetLogger(ctx).Panicf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	_log.Fatalf(format, args...)
}

func FatalfCtx(ctx context.Context, format string, args ...interface{}) {
	GetLogger(ctx).Fatalf(format, args...)
}
