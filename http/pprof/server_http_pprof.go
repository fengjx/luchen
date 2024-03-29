package pprof

import (
	"net/http/pprof"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/http/middleware"
)

const (
	defaultPrefix = "/debug/pprof"
)

type Handler struct {
	prefix string
	creds  map[string]string
}

func NewPprofHandler() Handler {
	return Handler{}
}

func (h Handler) Bind(router *luchen.HTTPServeMux) {
	h.routeRegister(router)
}

// Prefix url 前缀
func (h Handler) Prefix(prefix string) Handler {
	h.prefix = prefix
	return h
}

// BasicAuth basic 认证的用户名和密码，支持多组
func (h Handler) BasicAuth(creds map[string]string) Handler {
	h.creds = creds
	return h
}

func (h Handler) routeRegister(router *luchen.HTTPServeMux) {
	prefix := h.prefix
	if prefix == "" {
		prefix = defaultPrefix
	}
	mux := luchen.NewHTTPServeMux()
	mux.HandleFunc("/", pprof.Index)
	mux.HandleFunc("/cmdline", pprof.Cmdline)
	mux.HandleFunc("/profile", pprof.Profile)
	mux.HandleFunc("/symbol", pprof.Symbol)
	mux.HandleFunc("/trace", pprof.Trace)
	mux.Handle("/allocs", pprof.Handler("allocs"))
	mux.Handle("/block", pprof.Handler("block"))
	mux.Handle("/goroutine", pprof.Handler("goroutine"))
	mux.Handle("/heap", pprof.Handler("heap"))
	mux.Handle("/mutex", pprof.Handler("mutex"))
	mux.Handle("/threadcreate", pprof.Handler("threadcreate"))
	if len(h.creds) > 0 {
		mux.Use(middleware.BasicAuth("pprof", h.creds))
	}
	router.Sub(prefix, mux)
}
