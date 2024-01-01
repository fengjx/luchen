package luchen_test

import (
	"context"
	"net/http"

	"github.com/fengjx/go-halo/json"
	httptransport "github.com/go-kit/kit/transport/http"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/fengjx/luchen"
)

func newHelloHttpServer(serviceName, addr string) *luchen.HTTPServer {
	server := luchen.NewHTTPServer(
		serviceName,
		luchen.WithHTTPAddr(addr),
	).Handler(
		&helloHandler{},
	)
	return server
}

type helloHandler struct {
}

func (h *helloHandler) Bind(router luchen.HTTPRouter) {
	router.Handle("/say-hello", h.sayHello())
}

func (h *helloHandler) sayHello() *httptransport.Server {
	return httptransport.NewServer(
		makeSayHelloEndpoint(),
		luchen.DecodeJsonRequest[pb.HelloRequest],
		encodeSayHello,
	)
}

func encodeSayHello(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	jsonStr, err := json.ToJson(resp)
	if err != nil {
		return err
	}
	_, err = w.Write([]byte(jsonStr))
	return err
}
