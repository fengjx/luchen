package http

import (
	"github.com/go-chi/chi/v5"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/greetsvr/internal/hello"
	"github.com/fengjx/luchen/example/greetsvr/pb"
)

type greeterHandler struct {
}

func newGreeterHandler() *greeterHandler {
	return &greeterHandler{}
}

func (h *greeterHandler) Bind(router luchen.HTTPRouter) {
	router.Route("/hello", func(r chi.Router) {
		r.Handle("/say-hello", h.sayHello())
	})
}

func (h *greeterHandler) sayHello() *httptransport.Server {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	return httptransport.NewServer(
		hello.GetInst().Endpoints.MakeSayHelloEndpoint(),
		luchen.DecodeParamHTTPRequest[pb.HelloReq],
		luchen.CreateHTTPJSONEncoder(httpResponseWrapper),
		options...,
	)
}
