package luchen

import (
	"context"
	"time"
)

// Context 网关请求上下文
type Context struct {
	ctx context.Context
}

// NewContext 创建 context
func NewContext(ctx context.Context) *Context {
	return &Context{
		ctx: ctx,
	}
}

// WithValue 根据父context创建一个新的context，并设置新值
func WithValue(parent context.Context, key, val any) *Context {
	ctx := context.WithValue(parent, key, val)
	return &Context{ctx: ctx}
}

// Set 为 context 设置新值，在方法内部执行可以携带到方法外
// 参考测试用例：context_test#TestCtx
func (c *Context) Set(key, val any) {
	ctx := context.WithValue(c.ctx, key, val)
	c.ctx = ctx
}

// Deadline 同 context.Context 的 Deadline 方法
func (c *Context) Deadline() (deadline time.Time, ok bool) {
	return c.ctx.Deadline()
}

// Done 同 context.Context 的 Done 方法
func (c *Context) Done() <-chan struct{} {
	return c.ctx.Done()
}

// Err 同 context.Context 的 Err 方法
func (c *Context) Err() error {
	return c.ctx.Err()
}

// Value 根据key从上下文获取值
func (c *Context) Value(key any) any {
	return c.ctx.Value(key)
}
