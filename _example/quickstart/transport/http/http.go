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

func httpResponseWrapper(data interface{}) interface{} {
	res := &result{
		Msg:  "ok",
		Data: data,
	}
	return res
}

// 统一返回值处理
func encodeResponse(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res := &result{
		Msg:  "ok",
		Data: response,
	}
	log.InfoCtx(ctx, "http response", zap.Any("data", res))
	w.Header().Set("Content-Type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

// 统一异常处理
func errorEncoder(ctx context.Context, err error, w http.ResponseWriter) {
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
