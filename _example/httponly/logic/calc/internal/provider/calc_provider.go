package provider

import (
	"context"

	"github.com/fengjx/luchen/example/httponly/logic/calc/internal/service"
)

type CalcProvider struct {
}

func (p CalcProvider) Add(ctx context.Context, a int32, b int32) (int32, error) {
	return service.CalcSvc.Add(ctx, a, b)
}
