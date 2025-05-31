package gokit

import (
	"context"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

// App определяет основной интерфейс микросервиса для управления
// его жизненным циклом и функциональностью
type App interface {
	// Run запускает приложение и блокирует выполнение до получения
	// сигнала остановки или ошибки
	Run() error

	// Shutdown выполняет graceful shutdown приложения
	Shutdown(timeDuration time.Duration) error

	// Reconfigure перезагружает конфигурацию приложения без его остановки
	Reconfigure(timeDuraton time.Duration) error

	// RegisterGRPCServices регистрирует GRPC сервисы
	RegisterGRPCServices(_ grpc.ServiceRegistrar)

	// RegisterHandlersFromEndpoint регистрирует HTTP хендлеры для GRPC методов.
	// Используется для создания REST API поверх GRPC сервисов
	RegisterHandlersFromEndpoint(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) error

	// Health проверяет работоспособность критических компонентов
	Health(ctx context.Context) error
}
