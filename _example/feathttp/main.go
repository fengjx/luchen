package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/log"
)

// http server 功能演示

func main() {
	httpSvr := luchen.NewHTTPServer(
		luchen.WithServiceName("feathttp"),
		luchen.WithServerAddr(":8080"),
	).Use(
		timeMiddleware,
	).Handler(
		&helloHandler{},
	).Static("/assets/", "static")
	luchen.Start(httpSvr)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}

type helloHandler struct {
}

func (h *helloHandler) Bind(router *luchen.HTTPServeMux) {
	// curl http://localhost:8080/say-hello?name=fjx
	router.Handle("/say-hello", h.sayHello())

	// 为子路由添加中间件
	router.Sub("/log", func(sub *luchen.HTTPServeMux) {
		sub.Use(logMiddleware)
		sub.Handle("/say-hello", h.sayHello()) // curl http://localhost:8080/log/say-hello?name=fjx
	})
}

func (h *helloHandler) sayHello() *luchen.HTTPTransportServer {
	return luchen.NewHTTPTransportServer(
		makeSayHelloEndpoint(),
		decodeSayHello,
		encodeSayHello,
	)
}

func makeSayHelloEndpoint() luchen.Endpoint {
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

// 打印耗时中间件
func timeMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestURI := r.RequestURI
		start := time.Now()
		defer func() {
			log.Infof("take time: %s, %v", requestURI, time.Since(start))
		}()
		next.ServeHTTP(w, r)
	})
}

// 打印日志
func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.InfofCtx(r.Context(), "request %s", r.RequestURI)
		next.ServeHTTP(w, r)
	})
}
