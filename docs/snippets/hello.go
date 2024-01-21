func main() {
	// 启动一个 http server
	httpSvr := luchen.NewHTTPServer(
		"helloworld",
		luchen.WithHTTPAddr(":8080"),
	).Handler(
		&helloHandler{},	// 注册路由
	)
	luchen.Start(httpSvr)   // 启动服务

	// 监听系统 kill 信号
	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	// 优雅停机
	luchen.Stop()
}

type helloHandler struct {
}

// Bind 绑定 http 路由
func (h *helloHandler) Bind(router luchen.HTTPRouter) {
	router.Handle("/say-hello", h.sayHello())
}

// sayHello 绑定端点
func (h *helloHandler) sayHello() *httptransport.Server {
	return httptransport.NewServer(
		makeSayHelloEndpoint(),
		decodeSayHello,
		encodeSayHello,
	)
}

// makeSayHelloEndpoint 创建一个端点
func makeSayHelloEndpoint() kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		name := request.(string)
		response = "hello: " + name
		return
	}
}

// decodeSayHello 请求解码
func decodeSayHello(_ context.Context, r *http.Request) (interface{}, error) {
	name := r.URL.Query().Get("name")
	return name, nil
}

// encodeSayHello 响应参数编码
func encodeSayHello(_ context.Context, w http.ResponseWriter, resp interface{}) error {
	_, err := w.Write([]byte(resp.(string)))
	return err
}
