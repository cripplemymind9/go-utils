package server

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Server объединяет gRPC и Gateway серверы
type Server struct {
	grpc    *GRPCServer
	gateway *GatewayServer
}

// NewServer создает новый экземпляр сервера
func NewServer(grpcPort, httpPort int) *Server {
	grpcAddr := fmt.Sprintf("localhost:%d", grpcPort)
	return &Server{
		grpc:    NewGRPCServer(grpcPort),
		gateway: NewGatewayServer(httpPort, grpcAddr),
	}
}

// Start запускает оба сервера
func (s *Server) Start() error {
	// Запускаем gRPC сервер
	if err := s.grpc.Start(); err != nil {
		return fmt.Errorf("failed to start gRPC server: %w", err)
	}

	// Запускаем Gateway
	ctx := context.Background()
	if err := s.gateway.Start(ctx); err != nil {
		return fmt.Errorf("failed to start gateway: %w", err)
	}

	return nil
}

// Stop останавливает оба сервера
func (s *Server) Stop() error {
	ctx := context.Background()
	if err := s.gateway.Stop(ctx); err != nil {
		return fmt.Errorf("failed to stop gateway: %w", err)
	}
	s.grpc.Stop()
	return nil
}

// RunAndWait запускает серверы и ожидает сигнала для graceful shutdown
func (s *Server) RunAndWait() error {
	if err := s.Start(); err != nil {
		return err
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\nReceived shutdown signal")

	return s.Stop()
}
