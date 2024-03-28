package luchen

import (
	"net/http/pprof"
)

const (
	defaultPrefix = "/debug/pprof"
)

type PprofHandler struct {
	Prefix string
}

func (h PprofHandler) Bind(router *HTTPServeMux) {
	h.routeRegister(router)
}

func (h PprofHandler) routeRegister(router *HTTPServeMux) {
	prefix := h.Prefix
	if prefix == "" {
		prefix = defaultPrefix
	}
	router.HandleFunc(prefix+"/", pprof.Index)
	router.HandleFunc(prefix+"/cmdline", pprof.Cmdline)
	router.HandleFunc(prefix+"/profile", pprof.Index)
	router.HandleFunc(prefix+"/symbol", pprof.Symbol)
	router.HandleFunc(prefix+"/trace", pprof.Trace)
	router.Handle(prefix+"/allocs", pprof.Handler("allocs"))
	router.Handle(prefix+"/block", pprof.Handler("block"))
	router.Handle(prefix+"/goroutine", pprof.Handler("goroutine"))
	router.Handle(prefix+"/heap", pprof.Handler("heap"))
	router.Handle(prefix+"/mutex", pprof.Handler("mutex"))
	router.Handle(prefix+"/threadcreate", pprof.Handler("threadcreate"))
}
