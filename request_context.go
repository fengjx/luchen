package luchen

import (
	"context"
	"time"

	"github.com/fengjx/luchen/marshal"
)

type (
	requestHeaderKey struct{}

	httpRequestMarshallerCtxKey struct{}
)

type Header struct {
	Protocol    string
	Method      string
	Endpoint    string
	Host        string
	CLientIP    string
	TraceID     string
	StartTime   time.Time
	ContentType string
}

func withHeader(ctx context.Context, header *Header) context.Context {
	return context.WithValue(ctx, requestHeaderKey{}, header)
}

// GetHeader 返回当前请求的 header
func GetHeader(ctx context.Context) Header {
	return ctx.Value(requestHeaderKey{}).(Header)
}

func withMarshaller(ctx context.Context, marshaller marshal.Marshaller) context.Context {
	return context.WithValue(ctx, httpRequestMarshallerCtxKey{}, marshaller)
}

// Marshaller 返回当前请求的 marshaller
func Marshaller(ctx context.Context) marshal.Marshaller {
	return ctx.Value(httpRequestMarshallerCtxKey{}).(marshal.Marshaller)
}
