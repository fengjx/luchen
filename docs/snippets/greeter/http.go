type greeterHandler struct {
}

func newGreeterHandler() *greeterHandler {
	return &greeterHandler{}
}

func (h *greeterHandler) Bind(router luchen.HTTPRouter) {
	router.Route("/hello", func(r chi.Router) {
		r.Handle("/say-hello", h.sayHello())
	})
}

func (h *greeterHandler) sayHello() *httptransport.Server {
	options := []httptransport.ServerOption{
		httptransport.ServerErrorEncoder(errorEncoder),
	}
	return httptransport.NewServer(
		hello.GetInst().Endpoints.MakeSayHelloEndpoint(),
		luchen.DecodeParamHTTPRequest[pb.HelloReq],
		luchen.CreateHTTPJSONEncoder(httpResponseWrapper),
		options...,
	)
}
