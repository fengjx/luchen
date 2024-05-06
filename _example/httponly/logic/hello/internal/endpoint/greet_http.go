package endpoint

import (
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/httponly/logic/hello/internal/protocol"
	"github.com/fengjx/luchen/example/httponly/transport/http"
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
		luchen.DecodeHTTPParamRequest[protocol.HelloReq],
		luchen.EncodeHTTPJSONResponse(http.ResponseWrapper),
		options...,
	)
}
