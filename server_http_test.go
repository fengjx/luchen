package luchen_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fengjx/go-halo/json"
	httptransport "github.com/go-kit/kit/transport/http"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/fengjx/luchen"
)

func newHelloHttpServer(serviceName, addr string) *luchen.HTTPServer {
	server := luchen.NewHTTPServer(
		luchen.WithServiceName(serviceName),
		luchen.WithServerAddr(addr),
	).Handler(
		&helloHandler{},
	)
	return server
}

type helloHandler struct {
}

func (h *helloHandler) Bind(router *luchen.HTTPServeMux) {
	router.Handle("/say-hello", h.sayHello())
}

func (h *helloHandler) sayHello() *httptransport.Server {
	return luchen.NewHTTPHandler(
		makeSayHelloEndpoint(),
		luchen.DecodeHTTPJSONRequest[pb.HelloRequest],
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

func testStatic(t *testing.T) {
	httpServer := luchen.NewHTTPServer(
		luchen.WithServerAddr(":8080"),
	).Static("/static/", "docs/public")
	httpServer.Start()
}
