package luchen

import (
	"fmt"
)

var (
	// ErrSystem 系统错误
	ErrSystem         = &Errno{Code: 500, HTTPCode: 500, Msg: "系统错误"}
	ErrInvalidRequest = &Errno{Code: 400, HTTPCode: 400, Msg: "参数错误或类型不匹配"}
)

// Errno 错误编码定义
type Errno struct {
	Code     int    // 自定义错误码
	HTTPCode int    // http 错误码
	Msg      string // 错误信息
}

// Error 实现 error 接口
func (e *Errno) Error() string {
	return fmt.Sprintf("code:%d, httpcode:%d, msg:%s", e.Code, e.HTTPCode, e.Msg)
}
