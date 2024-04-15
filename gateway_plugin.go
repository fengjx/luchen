package luchen

import (
	"context"
	"net/http"
)

// GatewayPlugin 网关插件接口
type GatewayPlugin interface {
	// BeforeRoute 路由匹配前
	BeforeRoute(context.Context, *http.Request) (*http.Request, error)
	// AfterRoute 路由匹配后
	AfterRoute(context.Context, *http.Request) (*http.Request, error)
	// ModifyResponse 响应阶段处理
	ModifyResponse(context.Context, *http.Response) error
	// ErrorHandler 统一异常处理
	ErrorHandler(context.Context, http.ResponseWriter, *http.Request, error)
}

// UnimplementedGatewayPlugin 其他自定义插件如果不想实现所有接口，可以跟UnimplementedGatewayPlugin组合，只实现指定的方法即可
type UnimplementedGatewayPlugin struct {
}

// BeforeRoute nothing to do
func (p *UnimplementedGatewayPlugin) BeforeRoute(ctx context.Context, req *http.Request) (*http.Request, error) {
	return req, nil
}

// AfterRoute nothing to do
func (p *UnimplementedGatewayPlugin) AfterRoute(ctx context.Context, req *http.Request) (*http.Request, error) {
	return req, nil
}

// ModifyResponse nothing to do
func (p *UnimplementedGatewayPlugin) ModifyResponse(ctx context.Context, res *http.Response) error {
	return nil
}

// ErrorHandler nothing to do
func (p *UnimplementedGatewayPlugin) ErrorHandler(ctx context.Context, w http.ResponseWriter, req *http.Request, err error) {
	return
}

// GetContextErr 从上下文获取错误
func (p *UnimplementedGatewayPlugin) GetContextErr(ctx context.Context) error {
	if e, ok := ctx.Value(gatewayErrCtxKey{}).(error); ok {
		return e
	}
	return nil
}

// TraceGatewayPlugin 链路跟踪插件
type TraceGatewayPlugin struct {
	*UnimplementedGatewayPlugin
}

// BeforeRoute context 和 http header 注入 traceID
func (p *TraceGatewayPlugin) BeforeRoute(ctx context.Context, req *http.Request) (*http.Request, error) {
	r, traceID := TraceHTTPRequest(req)
	ctx = WithTraceID(ctx, traceID)
	return r, nil
}

// ModifyResponse 从上下文获取 tradeID 并写入 response header
func (p *TraceGatewayPlugin) ModifyResponse(ctx context.Context, res *http.Response) error {
	traceID := TraceID(ctx)
	res.Header.Set(TraceIDHeader, traceID)
	return nil
}
