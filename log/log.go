package log

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/fengjx/go-halo/halo"
	"github.com/fengjx/go-halo/json"
	"github.com/fengjx/go-halo/logger"
	"github.com/fengjx/go-halo/utils"
	kitlog "github.com/go-kit/log"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/env"
)

type (
	loggerKey struct{}
)

var (
	// LoggerCtxKey logger context key
	LoggerCtxKey = loggerKey{}

	_log    logger.Logger
	_logDir = filepath.Join("./", "logs")
)

func init() {
	level := logger.DebugLevel
	if env.IsProd() {
		level = logger.InfoLevel
	}
	logDir := os.Getenv("LUCHEN_LOG_DIR")
	if len(logDir) > 0 {
		_log = createFileLog(level, logDir)
		_log.SetLevel(level)
		return
	}
	if env.IsLocal() {
		_log = logger.NewConsole(zap.AddCallerSkip(2))
		_log.SetLevel(level)
		return
	}
	_log = createFileLog(level, GetLogDir())
	_log.SetLevel(level)
}

func createFileLog(level logger.Level, logDir string) logger.Logger {
	app := env.GetAppName()
	targetDir := filepath.Join(logDir, app)
	err := os.MkdirAll(targetDir, os.ModePerm)
	if err != nil {
		panic(err)
	}
	logfile := filepath.Join(targetDir, app+".log")
	appLog := logger.New(level, logfile, 1024, 7, zap.AddCallerSkip(2))
	log.Println("log file", logfile)
	appLog.Infof("log file: %s", logfile)
	return appLog
}

// GetLogDir 返回日志路径
func GetLogDir() string {
	logDir := os.Getenv("LUCHEN_LOG_DIR")
	if len(logDir) > 0 {
		return logDir
	}
	return _logDir
}

type kitLogger func(msg string, keysAndValues ...interface{})

// Log 日志打印实现
func (l kitLogger) Log(kv ...interface{}) error {
	fields := make(map[string]any)
	for i := 0; i < len(kv); i = i + 2 {
		k := kv[i]
		var v any
		n := i + 1
		if n <= len(kv) {
			v = kv[n]
		}
		fields[utils.ToString(k)] = v
	}
	jsonStr, _ := json.ToJson(fields)
	l("", jsonStr)
	return nil
}

// NewKitLogger returns a Go kit log.Logger that sends
func NewKitLogger(name string, level logger.Level) kitlog.Logger {
	targetLog := _log.With(zap.String("name", name))
	targetLog.SetLevel(level)
	var klog kitLogger
	switch level {
	case logger.DebugLevel:
		klog = targetLog.Debugf
	case logger.InfoLevel:
		klog = targetLog.Infof
	case logger.WarnLevel:
		klog = targetLog.Warnf
	case logger.ErrorLevel:
		klog = targetLog.Errorf
	case logger.DPanicLevel:
		klog = targetLog.Panicf
	case logger.PanicLevel:
		klog = targetLog.Panicf
	case logger.FatalLevel:
		klog = targetLog.Fatalf
	default:
		klog = targetLog.Infof
	}
	return klog
}

func GetLogger(ctx context.Context) logger.Logger {
	if lclog, ok := ctx.Value(LoggerCtxKey).(logger.Logger); ok {
		return lclog
	}
	goid := halo.GetGoID()
	lclog := _log.With(zap.Int64("goid", goid))
	return lclog
}

// WithLogger context 注入 logger
func WithLogger(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, LoggerCtxKey, GetLogger(ctx).With(fields...))
}
