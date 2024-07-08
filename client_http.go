package luchen

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/fengjx/go-halo/httpc"
	"github.com/samber/lo"
	"go.uber.org/zap"

	"github.com/fengjx/luchen/log"
)

var (
	retryHTTPCode = []int{
		0,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}

	httpClientCache     = make(map[string]*HTTPClient)
	httpClientCacheLock = newSegmentLock(10)
)

// HTTPClient 支持服务发现的 http client
type HTTPClient struct {
	serviceName string
	client      *httpc.Client
	selector    Selector
}

// HTTPResponse http 请求返回信息
type HTTPResponse = httpc.Response

// HTTPRequest http 请求参数
type HTTPRequest struct {
	Path   string
	Method string
	Header http.Header
	Params url.Values
	Form   url.Values
	Body   []byte
}

// GetHTTPClient 返回服务对应的 client
func GetHTTPClient(serviceName string) *HTTPClient {
	if cli, ok := httpClientCache[serviceName]; ok {
		return cli
	}
	lock := httpClientCacheLock.getLock(serviceName)
	lock.Lock()
	defer lock.Unlock()
	client := httpc.New(&httpc.Config{
		DefaultHeaders: map[string]string{
			"User-Agent": "luchen-http-client",
		},
		Timeout: defaultRequestTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: defaultConnectionTimeout,
			}).DialContext,
			MaxIdleConnsPerHost:   defaultMaxPoolSize,
			MaxIdleConns:          defaultPoolSize,
			IdleConnTimeout:       time.Second * 3,
			ExpectContinueTimeout: defaultConnectionTimeout,
		},
	})
	selector := GetEtcdV3Selector(serviceName)
	httpClient := &HTTPClient{
		serviceName: serviceName,
		client:      client,
		selector:    selector,
	}
	httpClientCache[serviceName] = httpClient
	return httpClient
}

// Call 发起 http 请求
func (c *HTTPClient) Call(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error) {
	var response *HTTPResponse
	var err error
	ch := make(chan error, 1)

	go func() {
		for i := 0; i < defaultRetries; i++ {
			response, err = c.call(ctx, req)
			if err != nil {
				ch <- nil
				break
			}
			if lo.Contains(retryHTTPCode, response.StatusCode()) {
				log.WarnCtx(ctx, "retry http call",
					zap.String("service_name", c.serviceName),
					zap.Any("req", req),
				)
				// retry
				continue
			}
		}
		ch <- nil
	}()

	select {
	case e := <-ch:
		err = e
	case <-ctx.Done():
		err = ctx.Err()
	}

	return response, err
}

func (c *HTTPClient) call(ctx context.Context, req *HTTPRequest) (*HTTPResponse, error) {
	node, err := c.selector.Next()
	if err != nil {
		return nil, err
	}
	rawurl := fmt.Sprintf("%s://%s%s", ProtocolHTTP, node.Addr, req.Path)
	var httpReq *http.Request
	var contentType string
	if req.Body != nil {
		contentType = "application/json"
		bodyBuf := bytes.NewBuffer(req.Body)
		httpReq, err = http.NewRequestWithContext(ctx, req.Method, rawurl, bodyBuf)
	} else if len(req.Form) > 0 {
		contentType = "application/x-www-form-urlencoded"
		httpReq, err = http.NewRequestWithContext(ctx, req.Method, rawurl, strings.NewReader(req.Form.Encode()))
	} else {
		httpReq, err = http.NewRequestWithContext(ctx, req.Method, rawurl, nil)
	}
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Content-Type", contentType)
	if len(req.Params) > 0 {
		httpReq.URL.RawQuery = req.Params.Encode()
	}
	httpReq, _ = TraceHTTPRequest(httpReq)
	return c.client.Request(httpReq)
}
