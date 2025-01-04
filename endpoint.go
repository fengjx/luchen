package luchen

import (
	"context"
	"fmt"
	"reflect"

	"github.com/go-kit/kit/endpoint"
	"google.golang.org/protobuf/proto"
)

type (
	// Endpoint alias for endpoint.Endpoint
	Endpoint = endpoint.Endpoint

	// Handler 服务接口处理器
	Handler func(ctx context.Context, req proto.Message) (resp proto.Message, err error)
)

// EdnpointDefine 端点定义信息
type EdnpointDefine struct {
	Name        string       // 端点名称
	Path        string       // http 请求路径
	ReqType     reflect.Type // 请求类型
	RspType     reflect.Type // 响应类型
	Handler     Handler      // 服务接口处理器
	Middlewares []Middleware // 中间件
}

// MakeEndpoint 包装 endpoint，添加类型安全和中间件支持
func MakeEndpoint(desc *EdnpointDefine) Endpoint {
	e := func(ctx context.Context, request interface{}) (response interface{}, err error) {
		// 使用 reflect 将 request 转换为 desc.ReqType 指定的类型
		reqVal := reflect.ValueOf(request)
		if reqVal.Type() != desc.ReqType {
			return nil, fmt.Errorf("%w: expected %v, got %T", ErrInvalidRequest, desc.ReqType, request)
		}
		return desc.Handler(ctx, request.(proto.Message))
	}
	return EndpointChain(e, desc.Middlewares...)
}

// EndpointChain Endpoint 中间件包装
func EndpointChain(e Endpoint, middlewares ...Middleware) Endpoint {
	if len(middlewares) == 0 {
		return e
	}
	return endpoint.Chain(middlewares[0], middlewares[1:]...)(e)
}
