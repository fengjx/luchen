package luchen_test

import (
	"context"
	"testing"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/fengjx/luchen"
)

func TestGRPCCall(t *testing.T) {
	serviceName := "rpc.hello"
	registrar := luchen.NewEtcdV3Registrar(
		newHelloGRPCServer(serviceName, ":0"),
	)
	registrar.Register()
	defer registrar.Deregister()

	grpcClient := luchen.GetGRPCClient(
		serviceName,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	greeterClient := pb.NewGreeterClient(grpcClient)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()
	helloResp, err := greeterClient.SayHello(ctx, &pb.HelloRequest{
		Name: "fengjx",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(helloResp.Message)
}
