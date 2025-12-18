package serve

import (
	"context"
	"time"

	"github.com/YumikoKawaii/shared/logger"
	"github.com/YumikoKawaii/shared/mysql"
	"github.com/YumikoKawaii/shared/redis"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	grpc_validator "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/validator"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	otel_sdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/connection_registry"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/receiver"
	"yumiko_kawaii.com/yine/applications/orchestrator/handlers/streamer"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/interceptor"
	"yumiko_kawaii.com/yine/applications/orchestrator/pkg/repository/uow"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
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

	ctx := context.Background()
	exporter, err := otlptracehttp.New(
		ctx,
		otlptracehttp.WithEndpoint("localhost:4318"),
		otlptracehttp.WithInsecure(),
	)
	if err != nil {
		panic(err)
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			semconv.ServiceName("receiver"),
		),
	)
	if err != nil {
		panic(err)
	}

	tp := otel_sdk.NewTracerProvider(
		otel_sdk.WithBatcher(
			exporter,
			otel_sdk.WithMaxExportBatchSize(512),
			otel_sdk.WithBatchTimeout(5*time.Second),
			otel_sdk.WithMaxQueueSize(2048),
		),
		// Service metadata
		otel_sdk.WithResource(res),

		// Sampling strategy
		otel_sdk.WithSampler(otel_sdk.TraceIDRatioBased(0.1)), // Sample 10%

		// Span limits
		otel_sdk.WithRawSpanLimits(otel_sdk.SpanLimits{
			AttributeCountLimit:         128,
			EventCountLimit:             128,
			LinkCountLimit:              128,
			AttributePerEventCountLimit: 128,
			AttributePerLinkCountLimit:  128,
		}),
	)
	otel.SetTracerProvider(tp)
	tracer := otel.Tracer("receiver_tracer")

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

	db := mysql.Initialize(&conf.MysqlCfg)
	redisCli := redis.Initialize(conf.RedisCfg)
	dbWorker := uow.New(db)
	connectionRegistry := connection_registry.NewRegistry(redisCli)
	messagePublisher := redis.NewPublisher(redisCli)
	srv := receiver.NewHandler(connectionRegistry, messagePublisher, dbWorker)

	if err = s.Register(
		srv,
	); err != nil {
		logger.Fatalf("error register server: %s", err.Error())
	}

	logger.Infof("http: %s", conf.Server.HTTP.String())
	logger.Infof("grpc: %s", conf.Server.GRPC.String())
	logger.Infof("serving...")

	if err = s.Serve(); err != nil {
		logger.Fatalf("error serving: %s", err.Error())
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
