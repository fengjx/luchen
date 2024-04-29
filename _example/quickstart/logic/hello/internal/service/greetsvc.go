package service

import (
	"context"
	"fmt"

	"github.com/fengjx/luchen/log"
	"go.uber.org/zap"
)

var GreetSvc *greetService

func init() {
	GreetSvc = &greetService{}
}

type greetService struct {
}

func (svc *greetService) SayHi(ctx context.Context, name string) (string, error) {
	log.InfoCtx(ctx, "say hi", zap.Any("name", name))
	return fmt.Sprintf("Hi: %s", name), nil
}
