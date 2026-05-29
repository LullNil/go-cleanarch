package auth

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/LullNil/go-cleanarch/config"
	"github.com/LullNil/go-cleanarch/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	healthpb "google.golang.org/grpc/health/grpc_health_v1"
)

const defaultRequestTimeout = 2 * time.Second

// Client calls the external auth service through gRPC.
type Client struct {
	conn           *grpc.ClientConn
	grpc           healthpb.HealthClient
	requestTimeout time.Duration
}

// NewGRPCClient creates a new auth gRPC integration client.
func NewGRPCClient(cfg config.AuthIntegration) (*Client, error) {
	const op = "integration.auth.NewGRPCClient"

	if strings.TrimSpace(cfg.GRPCTarget) == "" {
		return nil, fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	conn, err := grpc.NewClient(
		cfg.GRPCTarget,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: create grpc client: %w", op, err)
	}

	timeout := cfg.RequestTimeout
	if timeout <= 0 {
		timeout = defaultRequestTimeout
	}

	return &Client{
		conn:           conn,
		grpc:           healthpb.NewHealthClient(conn),
		requestTimeout: timeout,
	}, nil
}

// CanAccessEntity1 checks whether subjectID can access entity1ID.
func (c *Client) CanAccessEntity1(ctx context.Context, subjectID string, entity1ID int64) (bool, error) {
	const op = "integration.auth.CanAccessEntity1"

	if strings.TrimSpace(subjectID) == "" || entity1ID <= 0 {
		return false, fmt.Errorf("%s: %w", op, domain.ErrInvalidInput)
	}

	ctx, cancel := context.WithTimeout(ctx, c.requestTimeout)
	defer cancel()

	// In a real auth integration, replace this health check with the generated
	// authpb client call and map protobuf responses to the service port shape.
	resp, err := c.grpc.Check(ctx, &healthpb.HealthCheckRequest{
		Service: "auth.v1.AuthService",
	})
	if err != nil {
		return false, fmt.Errorf("%s: call auth service: %w", op, err)
	}

	return resp.GetStatus() == healthpb.HealthCheckResponse_SERVING, nil
}

// Close closes the underlying gRPC connection.
func (c *Client) Close() error {
	if c == nil || c.conn == nil {
		return nil
	}

	return c.conn.Close()
}
