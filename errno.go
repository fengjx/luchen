package luchen

import (
	"errors"
	"fmt"
	"net/http"
)

var (
	// ErrSystem 系统错误
	ErrSystem = NewErrnoWithHttpCode(http.StatusInternalServerError, "")
	// ErrNotImplemented 未实现
	ErrNotImplemented = NewErrnoWithHttpCode(http.StatusNotImplemented, "")
	// ErrBadGateway 网关错误
	ErrBadGateway = NewErrnoWithHttpCode(http.StatusBadGateway, "")
	// ErrServiceUnavailable 服务不可用
	ErrServiceUnavailable = NewErrnoWithHttpCode(http.StatusServiceUnavailable, "")
	// ErrGatewayTimeout 网关超时
	ErrGatewayTimeout = NewErrnoWithHttpCode(http.StatusGatewayTimeout, "")

	// ErrBadRequest 参数错误
	ErrBadRequest = NewErrnoWithHttpCode(http.StatusBadRequest, "")
	// ErrUnauthorized 未授权
	ErrUnauthorized = NewErrnoWithHttpCode(http.StatusUnauthorized, "")
	// ErrForbidden 没有权限访问
	ErrForbidden = NewErrnoWithHttpCode(http.StatusForbidden, "")
	// ErrNotFound 资源不存在
	ErrNotFound = NewErrnoWithHttpCode(http.StatusNotFound, "")
)

// NewErrno 创建错误编码
func NewErrno(httpcode int, code int, msg string) *Errno {
	ststusMsg := http.StatusText(httpcode)
	if msg == "" {
		msg = ststusMsg
	}
	return &Errno{HttpCode: httpcode, Code: code, Msg: msg}
}

// NewErrnoWithHttpCode 根据 http 状态码创建错误编码
func NewErrnoWithHttpCode(httpcode int, msg string) *Errno {
	return NewErrno(httpcode, httpcode, msg)
}

// Errno 错误编码定义
// Code: 0 - 500，使用 http 状态码规范，大于 1000 的错误码，表示自定义错误码
type Errno struct {
	HttpCode int    `json:"http_code"` // http 状态码
	Code     int    `json:"code"`      // 错误码
	Msg      string `json:"msg"`       // 错误信息
	Detail   string `json:"detail"`    // 详细错误信息
	Cause    error  `json:"-"`         // 原始错误，不序列化
}

// WithMsg 设置错误信息
func (e *Errno) WithMsg(msg string) *Errno {
	newErr := *e
	newErr.Msg = msg
	return &newErr
}

// WithDetail 设置详细错误信息
func (e *Errno) WithDetail(detail string) *Errno {
	newErr := *e
	newErr.Detail = detail
	return &newErr
}

// GetDetail 获取详细错误信息
func (e *Errno) GetDetail() string {
	if e.Detail != "" {
		return e.Detail
	}
	if e.Cause != nil {
		return fmt.Sprintf("%v", e.Cause)
	}
	return ""
}

// WithCause 设置原始错误
func (e *Errno) WithCause(err error) *Errno {
	newErr := *e
	newErr.Cause = err
	if e.Detail == "" && err != nil {
		newErr.Detail = err.Error()
	}
	return &newErr
}

// Unwrap 实现 errors.Unwrap 接口
func (e *Errno) Unwrap() error {
	return e.Cause
}

// Error 实现 error 接口
func (e *Errno) Error() string {
	if e.Detail != "" {
		return fmt.Sprintf("code:%d, msg:%s, detail:%s", e.Code, e.Msg, e.Detail)
	}
	return fmt.Sprintf("code:%d, msg:%s", e.Code, e.Msg)
}

// Format 实现 fmt.Formatter 接口，支持更丰富的格式化输出
func (e *Errno) Format(s fmt.State, verb rune) {
	switch verb {
	case 'v':
		if s.Flag('+') {
			// 详细格式
			fmt.Fprintf(s, "Errno{HttpCode: %d, Code: %d, Msg: %s, Detail: %s}",
				e.HttpCode, e.Code, e.Msg, e.Detail)
			if e.Cause != nil {
				fmt.Fprintf(s, "\nCaused by: %+v", e.Cause)
			}
			return
		}
		fallthrough
	case 's':
		fmt.Fprint(s, e.Error())
	}
}

// IsClientError 是否是客户端错误
func (e *Errno) IsClientError() bool {
	return e.HttpCode >= 400 && e.HttpCode < 500
}

// IsServerError 是否是服务端系统错误
func (e *Errno) IsServerError() bool {
	return e.HttpCode >= 500
}

// WrapError 包装普通错误为 Errno
func WrapError(err error) *Errno {
	if err == nil {
		return nil
	}
	var e *Errno
	if errors.As(err, &e) {
		return e
	}
	return ErrSystem.WithCause(err)
}

// FromError 从错误中提取 Errno
func FromError(err error) (*Errno, bool) {
	if err == nil {
		return nil, false
	}
	var errno *Errno
	if errors.As(err, &errno) {
		return errno, true
	}
	return nil, false
}
