package hello

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/httponly/logic/hello/internal/endpoint"
)

func Init(hs *luchen.HTTPServer) {
	endpoint.Init(hs)
}
