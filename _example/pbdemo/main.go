package main

import (
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/example/pbdemo/endpoint"
)

func main() {
	// 创建 http server
	hs := luchen.NewHTTPServer(
		luchen.WithServerAddr(":8080"),
	)

	// 注册 http 端点
	endpoint.RegisterGreeterHTTPHandler(hs)
	// 启动服务并监听 kill 信号
	hs.Start()
}
