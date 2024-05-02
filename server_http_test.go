package luchen_test

import (
	"context"
	"embed"
	"net/http"
	"testing"
	"time"

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
	return luchen.NewHTTPTransportServer(
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

//go:embed testdata
var embedFS embed.FS

func TestStatic(t *testing.T) {
	httpServer := luchen.NewHTTPServer(
		luchen.WithServerAddr(":8000"),
	).
		Static("/static/", "testdata/static").
		StaticFS("/fs/", luchen.Dir("testdata/static", true)).
		StaticFS("/", luchen.OnlyFilesFS(embedFS, false, "testdata/static"))

	if testing.Short() {
		select {
		case <-time.After(1 * time.Second * 10):
			httpServer.Stop()
		}
	}
	httpServer.Start()
}
