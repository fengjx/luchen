package service

import "github.com/fengjx/luchen/_example/greetsvr/service/hello"

func Init() {
	_ = hello.GetInst()
}
