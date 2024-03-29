package luchen

import (
	"net/http"
)

// HTTPMiddleware http 请求中间件
type HTTPMiddleware func(http.Handler) http.Handler

// HTTPServeMux http 路由
type HTTPServeMux struct {
	*http.ServeMux
	middlewares []HTTPMiddleware
	handler     http.Handler
}

// NewHTTPServeMux 创建一个 mux
func NewHTTPServeMux() *HTTPServeMux {
	mux := http.NewServeMux()
	router := &HTTPServeMux{
		ServeMux: mux,
	}
	router.then(mux)
	return router
}

// Use 注册中间件
func (mux *HTTPServeMux) Use(middlewares ...HTTPMiddleware) *HTTPServeMux {
	for _, middleware := range middlewares {
		mux.middlewares = append(mux.middlewares, middleware)
	}
	mux.then(mux.ServeMux)
	return mux
}

// Sub 注册子路由
func (mux *HTTPServeMux) Sub(prefix string, subMux *HTTPServeMux) {
	mux.Handle(prefix+"/", http.StripPrefix(prefix, subMux))
}

func (mux *HTTPServeMux) then(h http.Handler) {
	mux.handler = HandlerChain(h, mux.middlewares...)
}

func (mux *HTTPServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.handler.ServeHTTP(w, r)
}

// HandlerChain 使用中间件包装 handler
func HandlerChain(h http.Handler, middlewares ...HTTPMiddleware) http.Handler {
	size := len(middlewares)
	for i := range middlewares {
		h = middlewares[size-1-i](h)
	}
	return h
}
