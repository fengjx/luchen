package luchen

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"time"

	"github.com/fengjx/go-halo/httpc"
	"github.com/samber/lo"
	"go.uber.org/zap"
)

var (
	retryHTTPCode = []int{
		0,
		http.StatusNotImplemented,
		http.StatusBadGateway,
		http.StatusServiceUnavailable,
	}
)

type HTTPClient struct {
	serviceName string
	client      *httpc.Client
	selector    Selector
}

type HTTPResponse = httpc.Response

type HTTPRequest struct {
	Path   string
	Method string
	Header http.Header
	Params url.Values
	Form   url.Values
	Body   []byte
}

// GetHTTPClient 返回服务对应的 client
func GetHTTPClient(serviceName string, selector Selector) *HTTPClient {
	client := httpc.New(&httpc.Config{
		DefaultHeaders: map[string]string{
			"User-Agent": "luchen-http-client",
		},
		Timeout: time.Second * 3,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			DialContext: (&net.Dialer{
				Timeout: time.Second * 3,
			}).DialContext,
			MaxIdleConnsPerHost:   100,
			MaxIdleConns:          10,
			IdleConnTimeout:       time.Second * 3,
			ExpectContinueTimeout: time.Second,
		},
	})
	httpClient := &HTTPClient{
		serviceName: serviceName,
		client:      client,
		selector:    selector,
	}
	return httpClient
}

func (c *HTTPClient) Call(ctx context.Context, req *HTTPRequest) (response *HTTPResponse, err error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}
	for i := 0; i < 3; i++ {
		response, err = c.call(req)
		if err != nil {
			return nil, err
		}
		if lo.Contains(retryHTTPCode, response.StatusCode()) {
			RootLogger().Warn("retry http call",
				zap.String("service_name", c.serviceName),
				zap.Any("req", req),
			)
			// retry
			continue
		}
		return
	}
	return
}

func (c *HTTPClient) call(req *HTTPRequest) (*HTTPResponse, error) {
	node, err := c.selector.Next()
	if err != nil {
		return nil, err
	}
	rawurl := fmt.Sprintf("%s://%s%s", ProtocolHTTP, node.Addr, req.Path)
	var httpReq *http.Request
	if req.Body != nil {
		bodyBuf := bytes.NewBuffer(req.Body)
		httpReq, err = http.NewRequest(req.Method, rawurl, bodyBuf)
	} else {
		httpReq, err = http.NewRequest(req.Method, rawurl, nil)
	}
	if err != nil {
		return nil, err
	}
	if len(req.Form) > 0 {
		httpReq.Form = req.Form
	}
	if len(req.Params) > 0 {
		httpReq.URL.RawQuery = req.Params.Encode()
	}
	return c.client.Request(httpReq)
}
