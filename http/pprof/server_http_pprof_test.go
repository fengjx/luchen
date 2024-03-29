package pprof_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/http/pprof"
)

func TestPprofHandler(t *testing.T) {
	serviceName := "test_pprof_svr"
	server := luchen.NewHTTPServer(
		luchen.WithServiceName(serviceName),
		luchen.WithServerAddr(":8080"),
	).Handler(
		pprof.NewPprofHandler().BasicAuth(map[string]string{
			"foo": "bar",
		}),
	)
	if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Fatal(err)
	}

}
