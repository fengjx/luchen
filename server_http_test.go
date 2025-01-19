package luchen_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/fengjx/go-halo/json"
	"github.com/fengjx/luchen"
	"github.com/fengjx/luchen/env"
	"github.com/fengjx/luchen/types"
	"github.com/stretchr/testify/assert"
)

func TestHTTPServer(t *testing.T) {
	server := luchen.NewHTTPServer()
	server.Handle(&luchen.EndpointDefine{
		Path: "/test",
	})
}

func Test_errorEncoder(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantHttpCode  int
		wantErrCode   int32
		wantErrMsg    string
		wantHeaderSet bool
	}{
		{
			name:          "system error",
			err:           errors.New("system error"),
			wantHttpCode:  http.StatusInternalServerError,
			wantErrCode:   http.StatusInternalServerError,
			wantErrMsg:    http.StatusText(http.StatusInternalServerError),
			wantHeaderSet: true,
		},
		{
			name:          "custom error",
			err:           luchen.ErrBadRequest.WithMsg("invalid parameter"),
			wantHttpCode:  http.StatusBadRequest,
			wantErrCode:   http.StatusBadRequest,
			wantErrMsg:    "invalid parameter",
			wantHeaderSet: true,
		},
		{
			name:          "error with detail",
			err:           luchen.ErrSystem.WithDetail("database connection failed"),
			wantHttpCode:  http.StatusInternalServerError,
			wantErrCode:   http.StatusInternalServerError,
			wantErrMsg:    http.StatusText(http.StatusInternalServerError),
			wantHeaderSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx := context.Background()

			// 执行错误编码
			luchen.WriteError(ctx, w, tt.err)

			// 验证 HTTP 状态码
			assert.Equal(t, tt.wantHttpCode, w.Code)

			// 验证响应头
			rspMetaHeader := w.Header().Get(luchen.HeaderRspMeta)
			assert.Equal(t, tt.wantHeaderSet, rspMetaHeader != "")

			if tt.wantHeaderSet {
				// 解析响应元数据
				var rspMeta types.RspMeta
				err := json.FromJson(rspMetaHeader, &rspMeta)
				assert.NoError(t, err)

				// 验证错误码和错误信息
				assert.Equal(t, tt.wantErrCode, rspMeta.Code)
				assert.Equal(t, tt.wantErrMsg, rspMeta.Msg)

				// 验证其他必要字段
				assert.NotEmpty(t, rspMeta.ServerTime)
			}

			// 验证响应体为空
			assert.Empty(t, w.Body.String())
		})
	}
}

func TestWriteError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		wantHttpCode  int
		wantErrCode   int32
		wantErrMsg    string
		wantDetail    string
		wantHeaderSet bool
	}{
		{
			name:          "system error with cause",
			err:           luchen.ErrSystem.WithCause(errors.New("database error")),
			wantHttpCode:  http.StatusInternalServerError,
			wantErrCode:   http.StatusInternalServerError,
			wantErrMsg:    http.StatusText(http.StatusInternalServerError),
			wantDetail:    "database error",
			wantHeaderSet: true,
		},
		{
			name:          "bad request with detail",
			err:           luchen.ErrBadRequest.WithDetail("missing required field: name"),
			wantHttpCode:  http.StatusBadRequest,
			wantErrCode:   http.StatusBadRequest,
			wantErrMsg:    http.StatusText(http.StatusBadRequest),
			wantDetail:    "missing required field: name",
			wantHeaderSet: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			ctx := context.Background()

			// 执行错误写入
			luchen.WriteError(ctx, w, tt.err)

			// 验证 HTTP 状态码
			assert.Equal(t, tt.wantHttpCode, w.Code)

			// 验证响应头
			rspMetaHeader := w.Header().Get(luchen.HeaderRspMeta)
			assert.Equal(t, tt.wantHeaderSet, rspMetaHeader != "")

			if tt.wantHeaderSet {
				// 解析响应元数据
				var rspMeta types.RspMeta
				err := json.FromJson(rspMetaHeader, &rspMeta)
				assert.NoError(t, err)

				// 验证错误码和错误信息
				assert.Equal(t, tt.wantErrCode, rspMeta.Code)
				assert.Equal(t, tt.wantErrMsg, rspMeta.Msg)

				// 在非生产环境下验证详细错误信息
				if !env.IsProd() {
					assert.Equal(t, tt.wantDetail, rspMeta.Detail)
				}

				// 验证其他必要字段
				assert.NotEmpty(t, rspMeta.ServerTime)
				assert.NotEmpty(t, rspMeta.TraceId)
			}

			// 验证响应体为空
			assert.Empty(t, w.Body.String())
		})
	}
}
