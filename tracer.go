package luchen

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// TraceIDHeader traceID header key
	TraceIDHeader = "X-Trace-ID"
)

var (
	// TraceIDCtxKey traceID context key
	TraceIDCtxKey = struct{}{}
)

// TraceHTTPRequest 返回 traceID
// http 请求携带 traceID 处理
func TraceHTTPRequest(r *http.Request) (*http.Request, string) {
	traceID := r.Header.Get(TraceIDHeader)
	if traceID == "" {
		traceID = TraceID(r.Context())
	}
	if traceID == "" {
		traceID = uuid.NewString()
	}
	r.Header.Set(TraceIDHeader, traceID)
	ctx := WithTraceID(r.Context(), traceID)
	return r.WithContext(ctx), traceID
}

// TraceGRPC 返回 traceID
// grpc 请求携带 traceID 处理
func TraceGRPC(ctx context.Context, md metadata.MD) (context.Context, string) {
	traceID := uuid.NewString()
	if len(md.Get(TraceIDHeader)) > 0 {
		traceID = md.Get(TraceIDHeader)[0]
	}
	md.Set(TraceIDHeader, traceID)
	ctx = WithTraceID(ctx, traceID)
	return ctx, traceID
}

// TraceGRPCClient grpc client 携带 traceID
func TraceGRPCClient(ctx context.Context, method string, req, reply any, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
	traceID := TraceID(ctx)
	if traceID == "" {
		traceID = uuid.NewString()
		ctx = WithTraceID(ctx, traceID)
	}
	ctx = metadata.AppendToOutgoingContext(ctx, TraceIDHeader, traceID)
	return invoker(ctx, method, req, reply, cc, opts...)
}

// TraceID 从 context 获得 TraceID
func TraceID(ctx context.Context) string {
	value := ctx.Value(TraceIDCtxKey)
	if value == nil {
		return ""
	}
	return value.(string)
}

// WithTraceID context 注入 traceID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDCtxKey, traceID)
}
