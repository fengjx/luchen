package endpoint

import (
	"github.com/fengjx/luchen"
)

func Init(hs *luchen.HTTPServer) {
	// 注册 http 路由
	hs.Handler(
		&calcHandler{},
	)
}
