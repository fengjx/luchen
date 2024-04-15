package transport

import (
	"context"

	"github.com/go-kit/kit/sd"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/example/greetsvr/transport/grpc"
	"github.com/fengjx/luchen/example/greetsvr/transport/http"
)

var (
	registrar sd.Registrar
)

func Start(_ context.Context) {
	registrar = luchen.NewEtcdV3Registrar(
		grpc.GetServer(),
		http.GetServer(),
	)
	registrar.Register()
}

func Stop(_ context.Context) {
	registrar.Deregister()
}
