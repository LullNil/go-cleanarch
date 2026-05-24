package grpcserver

import (
	"context"
	"log/slog"
	"net"

	"github.com/LullNil/go-cleanarch/config"
	entity1grpc "github.com/LullNil/go-cleanarch/internal/delivery/grpc/entity1"
	"github.com/LullNil/go-cleanarch/internal/delivery/grpc/pb"

	"google.golang.org/grpc"
)

// Server wraps the gRPC server lifecycle.
type Server struct {
	log    *slog.Logger
	server *grpc.Server
	addr   string
}

// NewServer creates a new gRPC server.
func NewServer(cfg config.GRPCServer, log *slog.Logger, entity1Service entity1grpc.Service) *Server {
	server := grpc.NewServer(
		grpc.UnaryInterceptor(loggingUnaryInterceptor(log)),
	)

	pb.RegisterEntity1ServiceServer(server, entity1grpc.NewHandler(entity1Service))

	return &Server{
		log:    log,
		server: server,
		addr:   cfg.Port,
	}
}

// Run starts the gRPC server.
func (s *Server) Run() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	s.log.Info("starting grpc server", slog.String("addr", s.addr))
	return s.server.Serve(listener)
}

// Shutdown gracefully shuts down the gRPC server.
func (s *Server) Shutdown(ctx context.Context) error {
	done := make(chan struct{})
	go func() {
		s.server.GracefulStop()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		s.server.Stop()
		return ctx.Err()
	}
}
