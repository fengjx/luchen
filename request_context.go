package luchen

import (
	"context"
	"time"

	"github.com/fengjx/luchen/marshal"
)

type (
	requestEndpointCtxKey struct{}

	protocolCtxKey struct{}
	methodCtxKey   struct{}

	requestStartTimeCtxKey struct{}

	requestClientIP struct{}

	httpRequestMarshallerCtxKey struct{}
)

// RequestEndpoint 请求端点
func RequestEndpoint(ctx context.Context) string {
	val, ok := ctx.Value(requestEndpointCtxKey{}).(string)
	if !ok {
		return ""
	}
	return val
}

func withRequestEndpoint(ctx context.Context, action string) context.Context {
	return context.WithValue(ctx, requestEndpointCtxKey{}, action)
}

// RequestProtocol 请求协议
func RequestProtocol(ctx context.Context) string {
	val, ok := ctx.Value(protocolCtxKey{}).(string)
	if !ok {
		return ""
	}
	return val
}

func withRequestProtocol(ctx context.Context, protocol string) context.Context {
	return context.WithValue(ctx, protocolCtxKey{}, protocol)
}

// RequestMethod 请求方法
func RequestMethod(ctx context.Context) string {
	val, ok := ctx.Value(methodCtxKey{}).(string)
	if !ok {
		return ""
	}
	return val
}

func withMethod(ctx context.Context, method string) context.Context {
	return context.WithValue(ctx, methodCtxKey{}, method)
}

// RequestStartTime 请求开始时间
func RequestStartTime(ctx context.Context) time.Time {
	return ctx.Value(requestStartTimeCtxKey{}).(time.Time)
}

func withRequestStartTime(ctx context.Context, t time.Time) context.Context {
	return context.WithValue(ctx, requestStartTimeCtxKey{}, t)
}

func withRequestClientIP(ctx context.Context, ip string) context.Context {
	return context.WithValue(ctx, requestClientIP{}, ip)
}

// ClientIP 返回客户端IP
func ClientIP(ctx context.Context) string {
	return ctx.Value(requestClientIP{}).(string)
}

func withMarshaller(ctx context.Context, marshaller marshal.Marshaller) context.Context {
	return context.WithValue(ctx, httpRequestMarshallerCtxKey{}, marshaller)
}

// Marshaller 返回当前请求的 marshaller
func Marshaller(ctx context.Context) marshal.Marshaller {
	return ctx.Value(httpRequestMarshallerCtxKey{}).(marshal.Marshaller)
}
