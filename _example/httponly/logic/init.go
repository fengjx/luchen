package logic

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/httponly/logic/calc"
	"github.com/fengjx/luchen/example/httponly/logic/hello"
)

func Init(hs *luchen.HTTPServer) {
	hello.Init(hs)
	calc.Init(hs)
}
