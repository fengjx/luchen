package endpoint

import (
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/httponly/logic/calc/internal/protocol"
	"github.com/fengjx/luchen/example/httponly/transport/http"
)

type calcHandler struct {
}

func (h *calcHandler) Bind(router *luchen.HTTPServeMux) {
	router.Handle("/calc/add", h.sayHello())
}

func (h *calcHandler) sayHello() *httptransport.Server {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(http.ErrorEncoder),
	}
	return luchen.NewHTTPTransportServer(
		calcEdp.makeAddEndpoint(),
		luchen.DecodeHTTPParamRequest[*protocol.AddReq],
		luchen.EncodeHTTPJSONResponse(http.ResponseWrapper),
		options...,
	)
}
