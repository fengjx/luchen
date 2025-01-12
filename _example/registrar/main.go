package main

import (
	"time"

	"github.com/fengjx/go-halo/halo"
	"github.com/fengjx/luchen"

	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/example/registrar/endpoint"
)

func init() {
	if env.IsDev() {
		// 可以设置环境变量 LUCHEN_ETCD_ADDRESS 指定 etcd 地址
		env.SetDefaultEtcdAddress([]string{"192.168.6.121:2379"})
	}
}

func main() {

	// 创建 http server
	hs := luchen.NewHTTPServer(
		luchen.WithServiceName("http.helloworld"),
		luchen.WithServerAddr(":8080"),
	)

	// 创建 grpc server
	gs := luchen.NewGRPCServer(
		luchen.WithServiceName("grpc.helloworld"),
		luchen.WithServerAddr(":8088"),
	)

	// 注册 grpc 和 http 端点
	endpoint.RegisterGreeterGRPCHandler(gs)
	endpoint.RegisterGreeterHTTPHandler(hs)

	registrar := luchen.NewEtcdV3Registrar(
		hs,
		gs,
	)
	// 把服务注册到 etcd 并启动
	registrar.Register()
	// 阻塞服务并监听 kill 信号，收到 kill 信号后退出（最长等待10秒）
	halo.Wait(10 * time.Second)
}
