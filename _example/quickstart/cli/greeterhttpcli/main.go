package main

import (
	"context"
	"log"
	"net/http"
	"net/url"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/env"

	"github.com/fengjx/luchen/example/quickstart/connom/config"
)

func init() {
	if env.IsLocal() {
		// 可以设置环境变量 LUCHEN_ETCD_ADDRESS 指定 etcd 地址
		env.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func main() {
	sname := config.GetConfig().Server.HTTP.ServerName
	client := luchen.GetHTTPClient(sname)

	params := url.Values{}
	params.Set("name", "fengjx")
	req := &luchen.HTTPRequest{
		Path:   "/hello/say-hello",
		Method: http.MethodPost,
		Params: params,
	}
	response, err := client.Call(context.Background(), req)
	if err != nil {
		log.Fatal(err)
	}
	if !response.IsSuccess() {
		log.Fatal("http call not success")
	}
	log.Println(response.String())
}
