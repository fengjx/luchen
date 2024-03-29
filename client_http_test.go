package luchen_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/fengjx/go-halo/json"

	"github.com/fengjx/luchen"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func TestHTTPClient_Call(t *testing.T) {
	serviceName := "hello-test-http-call"
	registrar := luchen.NewEtcdV3Registrar(
		newHelloHttpServer(serviceName, ":0"),
	)
	registrar.Register()
	defer registrar.Deregister()
	client := luchen.GetHTTPClient(serviceName)
	body, _ := json.ToBytes(&pb.HelloRequest{
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
}
