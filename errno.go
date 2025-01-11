package luchen

import (
	"fmt"
	"net/http"
)

var (
	// ErrSystem 系统错误
	ErrSystem = NewErrno(http.StatusInternalServerError, "")
	// ErrBadRequest 参数错误
	ErrBadRequest = NewErrno(http.StatusBadRequest, "")
	// ErrForbidden 没有权限访问
	ErrForbidden = NewErrno(http.StatusForbidden, "")
)

// NewErrno 创建错误编码
func NewErrno(httpcode int, msg string) *Errno {
	ststusMsg := http.StatusText(httpcode)
	if msg == "" {
		msg = ststusMsg
	}
	return &Errno{Code: httpcode, Msg: msg}
}

// Errno 错误编码定义
// Code: 0 - 500，使用 http 状态码规范，大于 1000 的错误码，表示自定义错误码
type Errno struct {
	Code int    // 错误码
	Msg  string // 错误信息
}

// Error 实现 error 接口
func (e *Errno) Error() string {
	return fmt.Sprintf("code:%d, msg:%s", e.Code, e.Msg)
}

// IsServerError 是否是服务端系统错误
func (e *Errno) IsServerError() bool {
	return e.Code >= 500 && e.Code < 600
}
