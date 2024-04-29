package logic

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/logic/calc"
	"github.com/fengjx/luchen/example/quickstart/logic/hello"
)

func Init(hs *luchen.HTTPServer, gs *luchen.GRPCServer) {
	hello.Init(hs, gs)
	calc.Init(hs, gs)
}
