package serve

import (
	"context"

	"github.com/YumikoKawaii/shared/logger"
	"github.com/YumikoKawaii/shared/mysql"
	"github.com/YumikoKawaii/shared/redis"
	otel_tracer "github.com/YumikoKawaii/shared/tracer"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/connection_registry"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/receiver"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/streamer"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/interceptor"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/repository/uow"

	"github.com/spf13/cobra"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"yumiko_kawaii.com/yine/applications/orchestrator/config"
	"yumiko_kawaii.com/yine/applications/orchestrator/server"
)

func ServeReceiver(_ *cobra.Command, _ []string) {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Infof("Starting Receiver service initialization")

	ctx := context.Background()
	tracer, err := otel_tracer.Initialize(ctx, &conf.TracerConfig)
	logger.Infof("OpenTelemetry tracer initialized")

	traceInterceptor := interceptor.NewTracer(tracer)

	s := server.NewServer(conf.Server,
		grpc.KeepaliveParams(keepalive.ServerParameters{}),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
			traceInterceptor.Unary,
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	)

	logger.Infof("Initializing database and Redis connections")
	db := mysql.Initialize(&conf.MysqlCfg)
	redisCli, err := redis.Initialize(conf.RedisCfg)
	if err != nil {
		logger.Fatalf("error connecting redis: %s", err.Error())
	}
	dbWorker := uow.New(db)
	connectionRegistry := connection_registry.NewRegistry(redisCli)
	messagePublisher := redis.NewPublisher(redisCli)
	srv := receiver.NewHandler(connectionRegistry, messagePublisher, dbWorker)

	logger.Infof("Registering gRPC services")
	if err = s.Register(
		srv,
	); err != nil {
		logger.WithFields(logger.Fields{"error": err}).Fatalf("Error registering server")
	}

	logger.WithFields(logger.Fields{
		"http_addr": conf.Server.HTTP.Host,
		"http_port": conf.Server.HTTP.Port,
		"grpc_addr": conf.Server.GRPC.Host,
		"grpc_port": conf.Server.GRPC.Port,
	}).Infof("Starting Receiver server")

	if err = s.Serve(); err != nil {
		logger.WithFields(logger.Fields{"error": err}).Fatalf("Error starting server")
	}
}

func ServeStreamer(_ *cobra.Command, _ []string) {
	conf, err := config.Load()
	if err != nil {
		panic(err)
	}

	logger.Infof("Starting Streamer service initialization")

	s := server.NewServer(conf.Server,
		grpc.KeepaliveParams(keepalive.ServerParameters{}),
		grpc.ChainUnaryInterceptor(
			grpc_prometheus.UnaryServerInterceptor,
			grpc_validator.UnaryServerInterceptor(),
			grpc_recovery.UnaryServerInterceptor(),
		),
		grpc.ChainStreamInterceptor(
			grpc_prometheus.StreamServerInterceptor,
			grpc_validator.StreamServerInterceptor(),
			grpc_recovery.StreamServerInterceptor(),
		),
	)

	srv := streamer.NewHandler()

	logger.Infof("Registering gRPC services")
	if err = s.Register(
		srv,
	); err != nil {
		logger.WithFields(logger.Fields{"error": err}).Fatalf("Error registering servers")
	}

	logger.WithFields(logger.Fields{
		"grpc_addr": conf.Server.GRPC.Host,
		"grpc_port": conf.Server.GRPC.Port,
		"http_addr": conf.Server.HTTP.Host,
		"http_port": conf.Server.HTTP.Port,
	}).Infof("Starting Streamer server")

	if err = s.Serve(); err != nil {
		logger.WithFields(logger.Fields{"error": err}).Fatalf("Error starting server")
	}
}
