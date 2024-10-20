package luchen

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
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
	return func(next Endpoint) Endpoint {
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

// LogMiddleware 错误日志堆栈打印，放在第一个执行
func LogMiddleware() Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			resp, err := next(ctx, request)
			if err == nil {
				return resp, nil
			}
			var errn *Errno
			ok := errors.As(err, &errn)
			e := RequestEndpoint(ctx)
			if !ok {
				log.ErrorCtx(ctx, fmt.Sprintf("internal server Error: %+v", err), zap.Any("req", request), zap.String("endpoint", e))
			}
			return resp, err
		}
	}
}
