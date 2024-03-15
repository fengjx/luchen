package luchen

import (
	"context"
	"errors"
)

type AccessLogOpt struct {
	ContextFields map[string]any
	printResp     bool
	AccessLog     AccessLog
}

// AccessMiddleware 请求日志
func AccessMiddleware(opt AccessLogOpt) Middleware {
	return func(next Endpoint) Endpoint {
		return func(ctx context.Context, request interface{}) (response interface{}, err error) {
			fields := map[string]any{}
			for field, ctxKey := range opt.ContextFields {
				value := ctx.Value(ctxKey)
				fields[field] = value
			}
			fields["endpoint"] = RequestEndpoint(ctx)
			fields["protocol"] = RequestProtocol(ctx)
			fields["method"] = RequestMethod(ctx)
			fields["request"] = request

			response, err = next(ctx, request)
			if opt.printResp {
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

			accesslog := opt.AccessLog
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
