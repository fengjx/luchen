package http

import (
	"context"
	"errors"
	"net/http"

	"github.com/fengjx/go-halo/json"
	"github.com/fengjx/luchen/log"
	"go.uber.org/zap"

	"github.com/fengjx/luchen"
)

type result struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data any    `json:"data"`
}

// ResponseWrapper 响应数据包装
func ResponseWrapper(data interface{}) interface{} {
	res := &result{
		Msg:  "ok",
		Data: data,
	}
	return res
}

// ErrorEncoder 统一异常处理
func ErrorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
	log.ErrorCtx(ctx, "handler error", zap.Error(err))
	httpCode := 500
	msg := luchen.ErrSystem.Msg
	var errn *luchen.Errno
	ok := errors.As(err, &errn)
	if ok && errn.HTTPCode > 0 {
		httpCode = errn.HTTPCode
		msg = errn.Msg
	}
	w.WriteHeader(httpCode)
	res := &result{
		Code: httpCode,
		Msg:  msg,
	}
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		log.ErrorCtx(ctx, "write error msg fail", zap.Error(err))
	}
}
