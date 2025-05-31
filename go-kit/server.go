package gokit

import (
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	grpcServer   *grpc.Server
	httpServer   *http.Server
	gatewayMux   *runtime.ServeMux
	grpcPort     int
	httpPort     int
	grpcListener net.Listener
}

func NewServer(
	grpcServer *grpc.Server,
	cfg Config,
	muxOptions ...runtime.ServeMuxOption,
) (*Server, error) {
	grpcAddr := fmt.Sprintf("0.0.0.0:%d", cfg.GRPCPort)

	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		return nil, fmt.Errorf("failed to listen on %s: %w", grpcAddr, err)
	}

	gatewayMux := runtime.NewServeMux(muxOptions...)

	return &Server{
		grpcServer:   grpcServer,
		gatewayMux:   gatewayMux,
		grpcPort:     cfg.GRPCPort,
		httpPort:     cfg.HTTPPort,
		grpcListener: listener,
	}, nil
}

func (s *Server) GetMux() *runtime.ServeMux {
	return s.gatewayMux
}

func (s *Server) StartGRPCServer() {
	reflection.Register(s.grpcServer)

	go func() {
		if err := s.grpcServer.Serve(s.grpcListener); err != nil {
			log.Error().
				Err(err).
				Msg("Failed to serve gRPC")
		}
	}()
}

func (s *Server) StartGateway(app App) error {
	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.httpPort),
		Handler: s.gatewayMux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error().
				Err(err).
				Msg("Failed to serve HTTP gateway")
		}
	}()

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
	}

	if s.grpcServer != nil {
		s.grpcServer.GracefulStop()
	}

	return nil
}
