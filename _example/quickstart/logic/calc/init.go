package calc

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/quickstart/logic/calc/calcpub"
	"github.com/fengjx/luchen/example/quickstart/logic/calc/internal/endpoint"
	"github.com/fengjx/luchen/example/quickstart/logic/calc/internal/provider"
)

func Init(hs *luchen.HTTPServer, gs *luchen.GRPCServer) {
	calcpub.SetCalcAPI(provider.CalcProvider{})
	endpoint.Init(hs, gs)
}
