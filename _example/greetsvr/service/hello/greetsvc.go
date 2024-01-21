package hello

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	"github.com/fengjx/luchen"
)

type greetService struct {
}

func newGreetService() *greetService {
	return &greetService{}
}

func (svc *greetService) SayHi(ctx context.Context, name string) (string, error) {
	logger := luchen.Logger(ctx)
	logger.Info("say hi", zap.Any("name", name))
	return fmt.Sprintf("Hi: %s", name), nil
}
