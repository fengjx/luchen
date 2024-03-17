package luchen

import (
	"context"
	"errors"
	"time"
)

// 从 context 中获取值
type getCtxValueFn func(ctx context.Context) any

type AccessLogOpt struct {
	ContextFields map[string]getCtxValueFn
	PrintResp     bool
	AccessLog     AccessLog
}

// AccessMiddleware 请求日志
func AccessMiddleware(opt *AccessLogOpt) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			var accesslog AccessLog
			var contextFields map[string]getCtxValueFn
			var printResp bool
			if opt != nil {
				accesslog = opt.AccessLog
				contextFields = opt.ContextFields
				printResp = opt.PrintResp
			}
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
			if accesslog == nil {
				accesslog = NewAccessLog(1024, 7, 7)
			}
			accesslog.Print(fields)
			return
		}
	}
}

type AccessLog interface {
	Print(map[string]any)
}
