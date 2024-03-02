package luchen

import "net/http"

// HTTPMiddleware http 请求中间件
type HTTPMiddleware func(http.Handler) http.Handler

type ServeMux struct {
	*http.ServeMux
	middlewares []HTTPMiddleware
}

func NewServeMux() *ServeMux {
	mux := http.NewServeMux()
	return &ServeMux{
		ServeMux: mux,
	}
}

func (mux *ServeMux) Use(middlewares ...HTTPMiddleware) *ServeMux {
	for _, middleware := range middlewares {
		mux.middlewares = append(mux.middlewares, middleware)
	}
	return mux
}

func (mux *ServeMux) then(h http.Handler) http.Handler {
	size := len(mux.middlewares)
	for i := range mux.middlewares {
		h = mux.middlewares[size-1-i](h)
	}
	return h
}

func (mux *ServeMux) thenFunc(fn http.HandlerFunc) http.Handler {
	if fn == nil {
		return mux.then(nil)
	}
	return mux.then(fn)
}

func (mux *ServeMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mux.then(mux.ServeMux).ServeHTTP(w, r)
}
