package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GatewayRegistrar интерфейс для регистрации сервисов в gateway
type GatewayRegistrar interface {
	RegisterGateway(context.Context, *runtime.ServeMux, *grpc.ClientConn) error
}

// GatewayServer представляет собой HTTP сервер для gRPC-Gateway
type GatewayServer struct {
	httpServer *http.Server
	mux        *runtime.ServeMux
	port       int
	grpcAddr   string
	conn       *grpc.ClientConn
	registrars []GatewayRegistrar
}

// NewGatewayServer создает новый экземпляр GatewayServer
func NewGatewayServer(httpPort int, grpcAddr string) *GatewayServer {
	return &GatewayServer{
		mux:        runtime.NewServeMux(),
		port:       httpPort,
		grpcAddr:   grpcAddr,
		registrars: make([]GatewayRegistrar, 0),
	}
}

// RegisterService регистрирует новый сервис в gateway
func (s *GatewayServer) RegisterService(r GatewayRegistrar) {
	s.registrars = append(s.registrars, r)
}

// Start запускает HTTP сервер с gRPC-Gateway
func (s *GatewayServer) Start(ctx context.Context) error {
	conn, err := grpc.DialContext(
		ctx,
		s.grpcAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return fmt.Errorf("failed to dial gRPC server: %w", err)
	}
	s.conn = conn

	// Регистрируем все сервисы
	for _, registrar := range s.registrars {
		if err := registrar.RegisterGateway(ctx, s.mux, s.conn); err != nil {
			return fmt.Errorf("failed to register gateway: %w", err)
		}
	}

	s.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", s.port),
		Handler: s.mux,
	}

	go func() {
		if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("failed to serve HTTP: %v\n", err)
		}
	}()

	fmt.Printf("gRPC-Gateway started on port %d\n", s.port)
	return nil
}

// Stop останавливает HTTP сервер
func (s *GatewayServer) Stop(ctx context.Context) error {
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			return fmt.Errorf("failed to shutdown HTTP server: %w", err)
		}
		fmt.Println("gRPC-Gateway stopped gracefully")
	}
	if s.conn != nil {
		s.conn.Close()
	}
	return nil
}

// GetMux возвращает ServeMux для регистрации дополнительных handlers
func (s *GatewayServer) GetMux() *runtime.ServeMux {
	return s.mux
}
