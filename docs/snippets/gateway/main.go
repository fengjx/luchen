func main() {
	configFile, err := fs.Lookup("gateway.yaml", 3)
	if err != nil {
		luchen.RootLogger().Panic("config file not found", zap.Error(err))
	}
	config := luchen.MustLoadConfig[luchen.GatewayConfig](configFile)
	gateway := luchen.NewGateway(
		config,
		luchen.WithGatewayPlugin(
			&luchen.TraceGatewayPlugin{},
		),
	)
	luchen.Start(gateway)

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL)

	<-quit
	luchen.Stop()
}
