package luchen

import (
	"reflect"

	"github.com/go-kit/kit/endpoint"
)

type (
	// Endpoint alias for endpoint.Endpoint
	Endpoint = endpoint.Endpoint
)

// EndpointDefine 端点定义信息
type EndpointDefine struct {
	Name        string       // 端点名称
	Path        string       // http 请求路径
	ReqType     reflect.Type // 请求类型
	RspType     reflect.Type // 响应类型
	Endpoint    Endpoint     // 端点处理器
	Middlewares []Middleware // 中间件
}

// EndpointChain Endpoint 中间件包装
func EndpointChain(e Endpoint, middlewares ...Middleware) Endpoint {
	if len(middlewares) == 0 {
		return e
	}
	return endpoint.Chain(middlewares[0], middlewares[1:]...)(e)
}
