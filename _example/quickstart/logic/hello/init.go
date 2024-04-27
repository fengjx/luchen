package hello

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/logic/hello/internal/endpoint"
)

func Init(hs *luchen.HTTPServer, gs *luchen.GRPCServer) {
	endpoint.Init(hs, gs)
}
