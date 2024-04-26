package hello

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/logic/hello/endpoint"
)

func Init(hs *luchen.HTTPServer, gs *luchen.GRPCServer) {
	endpoint.Init(hs, gs)
}
