package luchen

import "context"

type (
	requestEndpointCtxKey struct{}

	protocolCtxKey struct{}
	methodCtxKey   struct{}
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
