package interceptor

import (
	"context"

	"github.com/YumikoKawaii/shared/logger"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpc_codes "google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

// metadataCarrier implements propagation.TextMapCarrier
// It adapts gRPC metadata.MD to work with OpenTelemetry's context propagation
type metadataCarrier struct {
	md metadata.MD
}

func (c *metadataCarrier) Get(key string) string {
	values := c.md.Get(key)
	if len(values) == 0 {
		return ""
	}
	return values[0]
}

func (c *metadataCarrier) Set(key, value string) {
	c.md.Set(key, value)
}

func (c *metadataCarrier) Keys() []string {
	keys := make([]string, 0, len(c.md))
	for k := range c.md {
		keys = append(keys, k)
	}
	return keys
}

type Tracer interface {
	Unary(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

func NewTracer(tracer trace.Tracer) Tracer {
	return &tracerImpl{
		tracer: tracer,
	}
}

type tracerImpl struct {
	tracer trace.Tracer
}

func (i *tracerImpl) Unary(ctx context.Context, request interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		ctx = otel.GetTextMapPropagator().Extract(ctx, &metadataCarrier{md: md})
	}

	ctx, span := i.tracer.Start(ctx, info.FullMethod,
		trace.WithSpanKind(trace.SpanKindServer),
	)
	defer span.End()

	span.SetAttributes(
		semconv.RPCSystemGRPC,
		semconv.RPCService(info.FullMethod),
	)

	spanCtx := span.SpanContext()
	traceID := spanCtx.TraceID().String()
	spanID := spanCtx.SpanID().String()
	logger.Infof("Handling request - TraceID: %s, SpanID: %s, Method: %s", traceID, spanID, info.FullMethod)

	resp, err := handler(ctx, request)

	// Record status
	if err != nil {
		span.RecordError(err)
		s, _ := status.FromError(err)
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(s.Code())))
		span.SetStatus(codes.Error, s.Message())
	} else {
		span.SetAttributes(semconv.RPCGRPCStatusCodeKey.Int(int(grpc_codes.OK)))
		span.SetStatus(codes.Ok, "")
	}

	return resp, err
}
