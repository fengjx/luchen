package service

import "github.com/fengjx/luchen/example/greetsvr/service/hello"

func Init() {
	_ = hello.GetInst()
}
