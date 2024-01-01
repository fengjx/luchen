package luchen

import "context"

type GRPCClient struct {
	serviceName string
	selector    Selector
}

func GetGRPCClient(serviceName string) *GRPCClient {
	grpcClient := &GRPCClient{
		serviceName: serviceName,
	}
	return grpcClient
}

func (c *GRPCClient) Call(ctx context.Context, method string, req any) (resp any, err error) {
	return
}
