package endpoint

import (
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/pb"
	"github.com/fengjx/luchen/example/quickstart/transport/http"
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
		luchen.DecodeHTTPParamRequest[*pb.AddReq],
		luchen.EncodeHTTPJSONResponse(http.ResponseWrapper),
		options...,
	)
}
