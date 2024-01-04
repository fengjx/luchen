package luchen

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

var (
	grpcClientCache     = make(map[string]*GRPCClient)
	grpcClientCacheLock = newSegmentLock(10)
)

type GRPCClient struct {
	selector    Selector
	pool        *pool
	dialOptions []grpc.DialOption
}

// GetGRPCClient 创建grpc客户端
func GetGRPCClient(serviceName string, opts ...grpc.DialOption) *GRPCClient {
	if cli, ok := grpcClientCache[serviceName]; ok {
		return cli
	}
	lock := grpcClientCacheLock.getLock(serviceName)
	lock.Lock()
	defer lock.Unlock()
	selector := GetEtcdV3Selector(serviceName)
	p := newPool(defaultPoolSize, defaultPoolTTL, defaultMaxPoolSize, defaultMaxPoolSize)
	client := &GRPCClient{
		selector:    selector,
		pool:        p,
		dialOptions: opts,
	}
	grpcClientCache[serviceName] = client
	return client
}

func (c *GRPCClient) Invoke(ctx context.Context, method string, args any, reply any, opts ...grpc.CallOption) error {
	conn, serviceInfo, err := c.getConn(ctx)
	if err != nil {
		return fmt.Errorf("error sending request: %v", err)
	}
	var grr error
	defer func() {
		c.pool.release(serviceInfo.Addr, conn, grr)
	}()

	ch := make(chan error, 1)

	go func() {
		err := conn.Invoke(ctx, method, args, reply, opts...)
		ch <- err
	}()

	select {
	case err := <-ch:
		grr = err
	case <-ctx.Done():
		grr = ctx.Err()
	}
	return grr
}

func (c *GRPCClient) NewStream(ctx context.Context, desc *grpc.StreamDesc, method string, opts ...grpc.CallOption) (grpc.ClientStream, error) {
	conn, serviceInfo, err := c.getConn(ctx)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %v", err)
	}
	var grr error
	defer func() {
		c.pool.release(serviceInfo.Addr, conn, grr)
	}()
	stream, grr := conn.NewStream(ctx, desc, method)
	return stream, grr
}

func (c *GRPCClient) next() (*ServiceInfo, error) {
	return c.selector.Next()
}

func (c *GRPCClient) getConn(ctx context.Context) (*poolConn, *ServiceInfo, error) {
	serviceInfo, err := c.next()
	if err != nil {
		return nil, nil, fmt.Errorf("find server node err: %v", err)
	}
	if serviceInfo == nil {
		return nil, nil, ErrNoServer
	}
	conn, err := c.pool.getConn(ctx, serviceInfo.Addr, c.dialOptions...)
	return conn, serviceInfo, err
}
