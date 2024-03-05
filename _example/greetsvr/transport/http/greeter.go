package http

import (
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/greetsvr/pb"
	"github.com/fengjx/luchen/example/greetsvr/service/hello"
)

type greeterHandler struct {
}

func newGreeterHandler() *greeterHandler {
	return &greeterHandler{}
}

func (h *greeterHandler) Bind(router *luchen.ServeMux) {
	router.Handle("/hello/say-hello", h.sayHello())
}

func (h *greeterHandler) sayHello() *httptransport.Server {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	return luchen.NewHTTPHandler(
		hello.GetInst().Endpoints.MakeSayHelloEndpoint(),
		luchen.DecodeHTTPParamRequest[pb.HelloReq],
		luchen.EncodeHTTPJSON(httpResponseWrapper),
		options...,
	)
}
