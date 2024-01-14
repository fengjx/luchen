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

func main() {
	httpSvr := luchen.NewHTTPServer(
		"helloworld",
		luchen.WithHTTPAddr(":8080"),
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

// Bind 绑定 http 请求路径
func (h *helloHandler) Bind(router luchen.HTTPRouter) {
	router.Handle("/say-hello", h.sayHello())
}

// sayHello 绑定端点
func (h *helloHandler) sayHello() *httptransport.Server {
	return httptransport.NewServer(
		makeSayHelloEndpoint(),
		decodeSayHello,
		encodeSayHello,
	)
}

// makeSayHelloEndpoint 创建一个端点
func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		name := request.(string)
		response = "hello: " + name
		return
	}
}

// decodeSayHello 请求解码
func decodeSayHello(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return name, nil
}

// encodeSayHello 响应参数编码
func encodeSayHello(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	_, err := w.Write([]byte(resp.(string)))
	return err
}
