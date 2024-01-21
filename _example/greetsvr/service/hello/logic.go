package hello

import (
	"context"
)

type helloLogic struct {
}

func newHelloLogic() *helloLogic {
	return &helloLogic{}
}

func (helloLogic) SayHello(ctx context.Context, name string) (string, error) {
	msg, err := GetInst().GreetSvc.SayHi(ctx, name)
	if err != nil {
		return "", err
	}
	return msg, nil
}
