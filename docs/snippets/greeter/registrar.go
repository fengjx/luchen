var (
	registrar sd.Registrar
)

func Start(_ context.Context) {
	registrar = luchen.NewEtcdV3Registrar(
		grpc.GetServer(),
		http.GetServer(),
	)
	registrar.Register()
}

func Stop(_ context.Context) {
	registrar.Deregister()
}
