package logic

import (
	"context"

	"github.com/fengjx/luchen/example/greetsvr/service/hello"
)

type helloLogic struct {
}

func newHelloLogic() *helloLogic {
	return &helloLogic{}
}

func (helloLogic) SayHello(ctx context.Context, name string) (string, error) {
	msg, err := hello.GetInst().GreetSvc.SayHi(ctx, name)
	if err != nil {
		return "", err
	}
	return msg, nil
}
