package main

import (
	"context"
	"log"
	"net/http"

	"github.com/fengjx/go-halo/json"
	"github.com/fengjx/luchen"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

func main() {
	client := luchen.GetHTTPClient("greeter")

	body, _ := json.ToBytes(&pb.HelloRequest{
		Name: "fengjx",
	})
	req := &luchen.HTTPRequest{
		Path:   "/hello/say-hello",
		Method: http.MethodPost,
		Body:   body,
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
