package luchen

import (
	"context"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

const (
	// TraceIDHeader traceID header key
	TraceIDHeader = "X-Trace-ID"
)

type (
	traceIDKey struct{}
)

var (
	// TraceIDCtxKey traceID context key
	TraceIDCtxKey = traceIDKey{}
)

// TraceHTTPRequest 返回 traceID
// http 请求携带 traceID 处理
func TraceHTTPRequest(r *http.Request) (*http.Request, string) {
	traceID := r.Header.Get(TraceIDHeader)
	if traceID == "" {
		traceID = TraceID(r.Context())
	}
	if traceID == "" {
		traceID = genTraceID()
	}
	r.Header.Set(TraceIDHeader, traceID)
	ctx := WithTraceID(r.Context(), traceID)
	return r.WithContext(ctx), traceID
}

// TraceGRPC 返回 traceID
// grpc 请求携带 traceID 处理
func TraceGRPC(ctx context.Context, md metadata.MD) (context.Context, string) {
	traceID := genTraceID()
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
		traceID = genTraceID()
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

// TraceIDOrNew 从 context 获得 TraceID，取不到则创建
func TraceIDOrNew(ctx context.Context) string {
	traceID := TraceID(ctx)
	if traceID == "" {
		return genTraceID()
	}
	return traceID
}

// WithTraceID context 注入 traceID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDCtxKey, traceID)
}

func genTraceID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
