package endpoint

import (
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/pb"
	"github.com/fengjx/luchen/example/quickstart/transport/http"
)

type greeterHandler struct {
}

func (h *greeterHandler) Bind(router *luchen.HTTPServeMux) {
	router.Handle("/hello/say-hello", h.sayHello())
}

func (h *greeterHandler) sayHello() *httptransport.Server {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(http.ErrorEncoder),
	}
	return luchen.NewHTTPTransportServer(
		greetEdp.makeSayHelloEndpoint(),
		luchen.DecodeHTTPParamRequest[pb.HelloReq],
		luchen.EncodeHTTPJSONResponse(http.ResponseWrapper),
		options...,
	)
}
