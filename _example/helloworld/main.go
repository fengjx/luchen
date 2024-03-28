package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	kitendpoint "github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"

	"github.com/fengjx/luchen"
)

// test cmd: curl -i http://localhost:8080/say-hello?name=luchen

func init() {
	if luchen.IsLocal() {
		luchen.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	httpSvr := luchen.NewHTTPServer(
		luchen.WithServiceName("helloworld"),
		luchen.WithServerAddr(":8080"),
	).Handler(
		&helloHandler{},
	)
	luchen.Start(httpSvr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}

type helloHandler struct {
}

func (h *helloHandler) Bind(router *luchen.HTTPServeMux) {
	router.Handle("/say-hello", h.sayHello())
}

func (h *helloHandler) sayHello() *httptransport.Server {
	return httptransport.NewServer(
		makeSayHelloEndpoint(),
		decodeSayHello,
		encodeSayHello,
	)
}

func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		name := request.(string)
		response = "hello: " + name
		return
	}
}

func decodeSayHello(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return name, nil
}

func encodeSayHello(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	_, err := w.Write([]byte(resp.(string)))
	return err
}
