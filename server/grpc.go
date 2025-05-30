package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GRPCServer представляет собой обертку над стандартным gRPC сервером
type GRPCServer struct {
	server   *grpc.Server
	port     int
	listener net.Listener
}

// NewGRPCServer создает новый экземпляр GRPCServer
func NewGRPCServer(port int, opts ...grpc.ServerOption) *GRPCServer {
	return &GRPCServer{
		server: grpc.NewServer(opts...),
		port:   port,
	}
}

// Start запускает gRPC сервер
func (s *GRPCServer) Start() error {
	addr := fmt.Sprintf(":%d", s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}
	s.listener = listener

	// Включаем reflection для удобства разработки
	reflection.Register(s.server)

	go func() {
		if err := s.server.Serve(listener); err != nil {
			fmt.Printf("failed to serve: %v\n", err)
		}
	}()

	fmt.Printf("gRPC server started on port %d\n", s.port)
	return nil
}

// Stop останавливает gRPC сервер
func (s *GRPCServer) Stop() {
	if s.server != nil {
		s.server.GracefulStop()
		fmt.Println("gRPC server stopped gracefully")
	}
}

// RunAndWait запускает сервер и ожидает сигнала для graceful shutdown
func (s *GRPCServer) RunAndWait() error {
	if err := s.Start(); err != nil {
		return err
	}

	// Настраиваем обработку сигналов для graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Ожидаем сигнал
	<-sigChan
	fmt.Println("\nReceived shutdown signal")

	s.Stop()
	return nil
}

// GetServer возвращает базовый gRPC сервер для регистрации сервисов
func (s *GRPCServer) GetServer() *grpc.Server {
	return s.server
}
