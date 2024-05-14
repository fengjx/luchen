package luchen

import (
	"context"
	"errors"
	"path"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"github.com/fengjx/luchen/log"
)

// GetValueFromContext 从 context 中获取值
type GetValueFromContext func(ctx context.Context) any

type AccessLogOpt struct {
	ContextFields map[string]GetValueFromContext
	PrintResp     bool
	AccessLog     AccessLog
	maxAge        int
}

// AccessMiddleware 请求日志
func AccessMiddleware(opt *AccessLogOpt) Middleware {
	return func(next Endpoint) Endpoint {
		var accesslog AccessLog
		var contextFields map[string]GetValueFromContext
		var printResp bool
		maxAge := 7
		if opt != nil {
			accesslog = opt.AccessLog
			contextFields = opt.ContextFields
			printResp = opt.PrintResp
			if opt.maxAge > 0 {
				maxAge = opt.maxAge
			}
		}
		if accesslog == nil {
			accesslog = NewAccessLog(1024, 7, maxAge)
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
	log *zap.Logger
}

func (impl accessLogImpl) Print(fields map[string]any) {
	var zf []zap.Field
	for field, value := range fields {
		zf = append(zf, zap.Any(field, value))
	}
	impl.log.Info("", zf...)
}

// NewAccessLog 创建一个 AccessLog
func NewAccessLog(maxSizeMB int, maxBackups int, maxAge int) AccessLog {
	logPath := path.Join(log.GetLogDir(), "access.log")
	w := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logPath,
		MaxSize:    maxSizeMB,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
	})
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "time"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")
	encoderConfig.FunctionKey = ""
	encoderConfig.LevelKey = ""
	encoderConfig.MessageKey = ""
	encoderConfig.NameKey = ""
	encoderConfig.CallerKey = ""
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		w,
		zapcore.InfoLevel,
	)
	l := zap.New(core, zap.AddCaller())
	return &accessLogImpl{log: l}
}
