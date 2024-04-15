package pprof

import (
	"expvar"
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

func NewHandler() Handler {
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
	router.Sub(prefix, func(sub *luchen.HTTPServeMux) {
		sub.HandleFunc("/", pprof.Index)
		sub.HandleFunc("/cmdline", pprof.Cmdline)
		sub.HandleFunc("/profile", pprof.Profile)
		sub.HandleFunc("/symbol", pprof.Symbol)
		sub.HandleFunc("/trace", pprof.Trace)

		sub.Handle("/vars", expvar.Handler())
		sub.Handle("/allocs", pprof.Handler("allocs"))
		sub.Handle("/block", pprof.Handler("block"))
		sub.Handle("/goroutine", pprof.Handler("goroutine"))
		sub.Handle("/heap", pprof.Handler("heap"))
		sub.Handle("/mutex", pprof.Handler("mutex"))
		sub.Handle("/threadcreate", pprof.Handler("threadcreate"))
		if len(h.creds) > 0 {
			sub.Use(middleware.BasicAuth("pprof", h.creds))
		}
	})
}
