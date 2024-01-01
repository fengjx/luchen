package luchen

import (
	"fmt"
)

var (
	SystemErr = &Errno{Code: 500, HttpCode: 500, Msg: "系统错误"}
)

type Errno struct {
	Code     int
	HttpCode int
	Msg      string
}

func (e *Errno) Error() string {
	return fmt.Sprintf("code:%d, httpcode:%d, msg:%s", e.Code, e.HttpCode, e.Msg)
}
