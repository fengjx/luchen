package luchen_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fengjx/go-halo/json"

	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/_example/greetsvr/pb"
)

func TestHTTPClient_Call(t *testing.T) {
	serviceName := "hello"
	registrar := luchen.NewEtcdV3Registrar(
		newHelloHttpServer(serviceName, ":0"),
	)
	registrar.Register()
	selector := luchen.NewEtcdV3Selector(serviceName)
	client := luchen.GetHTTPClient(serviceName, selector)
	body, _ := json.ToBytes(&pb.HelloReq{
		Name: "fengjx",
	})
	req := &luchen.HTTPRequest{
		Path:   "/say-hello",
		Method: http.MethodPost,
		Body:   body,
	}
	response, err := client.Call(context.Background(), req)
	if err != nil {
		t.Fatal(err)
	}
	if !response.IsSuccess() {
		t.Fatal("http call not success")
	}
	t.Log(response.String())
	registrar.Deregister()
}
