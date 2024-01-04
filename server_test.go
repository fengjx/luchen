package luchen_test

import (
	"context"

	kitendpoint "github.com/go-kit/kit/endpoint"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"

	"github.com/fengjx/luchen"
)

func init() {
	if luchen.IsLocal() {
		luchen.SetDefaultEtcdAddress([]string{"host.etcd.dev:2379"})
	}
}

func startTestServer() {
	httpSvr := newHelloHttpServer("hello", ":0")
	grpcSvr := newHelloGRPCServer("rpc.hello", ":0")
	luchen.Start(httpSvr, grpcSvr)
}

func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		helloReq := request.(*pb.HelloRequest)
		helloReply := &pb.HelloReply{
			Message: "Hi: " + helloReq.Name,
		}
		return helloReply, nil
	}
}
