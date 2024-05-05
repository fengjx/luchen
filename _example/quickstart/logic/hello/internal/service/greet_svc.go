package service

import (
	"context"
	"fmt"

	"github.com/fengjx/luchen/log"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/example/quickstart/logic/calc/calcpub"
)

var GreetSvc *greetService

func init() {
	GreetSvc = &greetService{}
}

type greetService struct {
	count int32
}

func (svc *greetService) SayHi(ctx context.Context, name string) (string, error) {
	log.InfoCtx(ctx, "say hi",
		zap.String("name", name),
		zap.Int32("count", svc.count),
	)
	count, err := calcpub.CalcAPI.Add(ctx, svc.count, 1)
	if err != nil {
		return "", err
	}
	svc.count = count
	return fmt.Sprintf("Hi: %s, %d", name, svc.count), nil
}
