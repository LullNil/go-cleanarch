package grpcserver

import (
	"context"
	"log/slog"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func loggingUnaryInterceptor(log *slog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		resp, err := handler(ctx, req)
		if err != nil {
			log.Error("grpc request failed",
				"method", info.FullMethod,
				"code", status.Code(err).String(),
				"error", err,
				"request_id", metadataValue(ctx, "x-request-id"),
				"trace_id", traceID(ctx),
			)
		}

		return resp, err
	}
}

func metadataValue(ctx context.Context, key string) string {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return ""
	}

	values := md.Get(key)
	if len(values) == 0 {
		return ""
	}

	return values[0]
}

func traceID(ctx context.Context) string {
	if traceID := metadataValue(ctx, "x-trace-id"); traceID != "" {
		return traceID
	}

	traceparent := metadataValue(ctx, "traceparent")
	if len(traceparent) >= 55 {
		return traceparent[3:35]
	}

	return ""
}
