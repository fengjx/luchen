package main

import (
	"context"
	"reflect"

	"github.com/fengjx/luchen"
)

func main() {
	// 创建 http server
	hs := luchen.NewHTTPServer(
		luchen.WithServerAddr(":8080"),
	)
	def := &luchen.EndpointDefine{
		Endpoint: sayHello,
		Path:     "/say-hello",
		ReqType:  reflect.TypeOf(&sayHelloReq{}),
		RspType:  reflect.TypeOf(&sayHelloRsp{}),
	}
	hs.Handle(def)
	// 启动服务并监听 kill 信号
	hs.Start()
}

func sayHello(ctx context.Context, request any) (response any, err error) {
	req := request.(*sayHelloReq)
	response = &sayHelloRsp{
		Msg: "hello " + req.Name,
	}
	return
}

type sayHelloReq struct {
	Name string
}

type sayHelloRsp struct {
	Msg string
}
