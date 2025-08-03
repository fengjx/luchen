package luchen

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/fengjx/go-halo/errs"
	"github.com/go-kit/kit/endpoint"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/log"
)

// Middleware alias for endpoint.Endpoint
type Middleware = endpoint.Middleware

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
			h := GetHeader(ctx)
			fields["endpoint"] = h.Endpoint
			fields["protocol"] = h.Protocol
			fields["ip"] = h.CLientIP
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
			startTime := h.StartTime
			fields["rt"] = time.Since(startTime).Nanoseconds()
			fields["rts"] = time.Since(startTime).String()
			accesslog.Print(fields)
			return
		}
	}
}

// LogMiddleware 错误日志堆栈打印，放在第一个执行
func LogMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		resp, err := next(ctx, request)
		if err == nil {
			return resp, nil
		}
		var errn *Errno
		ok := errors.As(err, &errn)
		h := GetHeader(ctx)
		if !ok || errn.IsServerError() {
			log.ErrorCtx(ctx,
				fmt.Sprintf("internal server Error: %+v", err),
				zap.Any("req", request),
				zap.String("endpoint", h.Endpoint),
				zap.Stack("stack"),
			)
		}
		return resp, err
	}
}

// RecoverMiddleware panic 处理
func RecoverMiddleware(next endpoint.Endpoint) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (resp interface{}, err error) {
		defer errs.RecoverFunc(func(perr any, stack *errs.Stack) {
			log.ErrorCtx(ctx, "server panic", zap.Any("req", request), zap.Any("error", perr), zap.String("stack", fmt.Sprintf("%+v", stack)))
			os.Stderr.Write([]byte(fmt.Sprintf("server panic: %v - stack: %+v", perr, stack)))
			resp = nil
			if env.IsProd() {
				err = ErrSystem
			} else {
				err = ErrSystem.WithDetail(fmt.Sprintf("%v", stack))
			}
		})
		resp, err = next(ctx, request)
		return resp, err
	}
}
