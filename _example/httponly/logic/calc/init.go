package calc

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/httponly/logic/calc/calcpub"
	"github.com/fengjx/luchen/example/httponly/logic/calc/internal/endpoint"
	"github.com/fengjx/luchen/example/httponly/logic/calc/internal/provider"
)

func Init(hs *luchen.HTTPServer) {
	calcpub.SetCalcAPI(provider.CalcProvider{})
	endpoint.Init(hs)
}
