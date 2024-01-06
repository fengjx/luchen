package luchen

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
)

const (
	// TraceIDHeader traceID header key
	TraceIDHeader = "X-Trace-ID"
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
	if len(md.Get(string(TraceIDCtxKey))) > 0 {
		traceID = md.Get(TraceIDHeader)[0]
	}
	md.Set(TraceIDHeader, traceID)
	ctx = WithTraceID(ctx, traceID)
	return ctx, traceID
}
