package luchen

import (
	"path/filepath"

	"github.com/fengjx/go-halo/logger"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/log"
)

type AccessLog interface {
	Print(map[string]any)
}

type accessLogImpl struct {
	log logger.Logger
}

func (impl accessLogImpl) Print(fields map[string]any) {
	var zf []zap.Field
	for field, value := range fields {
		zf = append(zf, zap.Any(field, value))
	}
	impl.log.Info("", zf...)
}

// NewAccessLog 创建一个 AccessLog
func NewAccessLog(maxSizeMB int, maxBackups int, maxDay int) AccessLog {
	logFile := filepath.Join(log.GetLogDir(), "access.log")
	l := logger.New(&logger.Options{
		LogFile:    logFile,
		MaxSizeMB:  maxSizeMB,
		MaxBackups: maxBackups,
		MaxDays:    maxDay,
		Thin:       true,
	})
	return &accessLogImpl{log: l}
}
