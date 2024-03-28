package luchen_test

import (
	"errors"
	"net/http"
	"testing"

	"github.com/fengjx/luchen"
)

func TestPprofHandler(t *testing.T) {
	serviceName := "test_pprof_svr"
	server := newHelloHttpServer(serviceName, ":8080")
	server.Handler(&luchen.PprofHandler{
		Prefix: "/debug/pprof",
	})
	if err := server.Start(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		t.Fatal(err)
	}

}
