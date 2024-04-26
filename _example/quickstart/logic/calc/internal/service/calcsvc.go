package service

import (
	"context"
)

var CalcSvc *calcService

func init() {
	CalcSvc = &calcService{}
}

type calcService struct {
}

func (svc *calcService) Add(_ context.Context, a, b int) (int, error) {
	return a + b, nil
}
