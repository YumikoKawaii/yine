package serve

import (
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/receiver"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/streamer"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"yumiko_kawaii.com/yine/applications/orchestrator/config"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/logger"
	"yumiko_kawaii.com/yine/applications/orchestrator/server"
)

func ServeReceiver(_ *cobra.Command, _ []string) {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	l, err := logger.NewLogger(conf.Logger)
	if err != nil {
		logger.WithFields(logger.Fields{"error": err}).Fatalf("Error initializing logger")
	}

	zapLogger := l.GetDelegate().(*zap.SugaredLogger).Desugar()

	recoveryOpt := grpc_recovery.WithRecoveryHandler(recoveryHandler(zapLogger))

	s := server.NewServer(conf.Server,
		grpc.KeepaliveParams(keepalive.ServerParameters{}),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(recoveryOpt),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(recoveryOpt),
		),
	)

	srv := receiver.NewHandler()

	if err = s.Register(
		srv,
	); err != nil {
		l.WithFields(logger.Fields{"error": err}).Fatalf("Error register servers")
	}

	l.WithFields(logger.Fields{"grpc_addr": conf.Server.GRPC.Host}).
		WithFields(logger.Fields{"grpc_port": conf.Server.GRPC.Port}).
		WithFields(logger.Fields{"http_addr": conf.Server.HTTP.Host}).
		WithFields(logger.Fields{"http_port": conf.Server.HTTP.Port}).
		Infof("Starting server...")

	if err = s.Serve(); err != nil {
		l.WithFields(logger.Fields{"error": err}).Fatalf("Error starting server")
	}
}

func ServeStreamer(_ *cobra.Command, _ []string) {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	l, err := logger.NewLogger(conf.Logger)
	if err != nil {
		logger.WithFields(logger.Fields{"error": err}).Fatalf("Error initializing logger")
	}

	zapLogger := l.GetDelegate().(*zap.SugaredLogger).Desugar()

	recoveryOpt := grpc_recovery.WithRecoveryHandler(recoveryHandler(zapLogger))

	s := server.NewServer(conf.Server,
		grpc.KeepaliveParams(keepalive.ServerParameters{}),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(recoveryOpt),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(recoveryOpt),
		),
	)

	srv := streamer.NewHandler()

	if err = s.Register(
		srv,
	); err != nil {
		l.WithFields(logger.Fields{"error": err}).Fatalf("Error register servers")
	}

	l.WithFields(logger.Fields{"grpc_addr": conf.Server.GRPC.Host}).
		WithFields(logger.Fields{"grpc_port": conf.Server.GRPC.Port}).
		WithFields(logger.Fields{"http_addr": conf.Server.HTTP.Host}).
		WithFields(logger.Fields{"http_port": conf.Server.HTTP.Port}).
		Infof("Starting server...")

	if err = s.Serve(); err != nil {
		l.WithFields(logger.Fields{"error": err}).Fatalf("Error starting server")
	}
}
