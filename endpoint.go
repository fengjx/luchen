package luchen

import (
	"context"
	"errors"
	"path/filepath"
	"time"

	"github.com/fengjx/go-halo/logger"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/log"
)

// GetValueFromContext 从 context 中获取值
type GetValueFromContext func(ctx context.Context) any

type AccessLogOpt struct {
	ContextFields map[string]GetValueFromContext
	PrintResp     bool
	AccessLog     AccessLog
	MaxDay        int
}

// AccessMiddleware 请求日志
func AccessMiddleware(opt *AccessLogOpt) Middleware {
	return func(next Endpoint) Endpoint {
		var accesslog AccessLog
		var contextFields map[string]GetValueFromContext
		var printResp bool
		maxDay := 7
		if opt != nil {
			accesslog = opt.AccessLog
			contextFields = opt.ContextFields
			printResp = opt.PrintResp
			if opt.MaxDay > 0 {
				maxDay = opt.MaxDay
			}
		}
		if accesslog == nil {
			accesslog = NewAccessLog(10*1024, maxDay, maxDay)
		}
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			fields := map[string]any{}
			for field, fn := range contextFields {
				value := fn(ctx)
				fields[field] = value
			}
			fields["endpoint"] = RequestEndpoint(ctx)
			fields["protocol"] = RequestProtocol(ctx)
			fields["method"] = RequestMethod(ctx)
			fields["ip"] = ClientIP(ctx)
			fields["request"] = request

			response, err = next(ctx, request)
			if printResp {
				fields["response"] = response
			}
			code := 0
			if err != nil {
				var errn *Errno
				ok := errors.As(err, &errn)
				if ok {
					code = errn.Code
				}
				fields["err"] = err.Error()
			}
			fields["code"] = code
			startTime := RequestStartTime(ctx)
			fields["rt"] = time.Since(startTime).Nanoseconds()
			fields["rts"] = time.Since(startTime).String()
			accesslog.Print(fields)
			return
		}
	}
}

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
