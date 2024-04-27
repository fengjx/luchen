package calc

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/logic/calc/internal/endpoint"
)

func Init(hs *luchen.HTTPServer, gs *luchen.GRPCServer) {
	endpoint.Init(hs, gs)
}
