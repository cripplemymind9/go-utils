# go-utils

## Go-Kit: Фреймворк для микросервисов

Легковесный фреймворк для создания gRPC и REST микросервисов.

### Быстрый старт

```go
package main

import (
	"github.com/cripplemymind9/go-utils/go-kit"
)

func main() {
	// Создаем новый Runner
	runner := gokit.NewRunner()
	
	// Создаем реализацию App
	app := NewMyApp()
	
	// Запускаем сервис
	if err := runner.Run(app); err != nil {
		panic(err)
	}
}
```

### Реализация интерфейса App

```go
// Реализуем интерфейс App
type MyApp struct {}

// Run запускает приложение
func (a *MyApp) Run() error {
	// Логика запуска приложения
	return nil
}

// RegisterGRPCServices регистрирует gRPC сервисы
func (a *MyApp) RegisterGRPCServices(server grpc.ServiceRegistrar) {
	// Регистрируем gRPC сервисы
	pb.RegisterMyServiceServer(server, a)
}

// RegisterHandlersFromEndpoint регистрирует REST эндпоинты для gRPC сервисов
func (a *MyApp) RegisterHandlersFromEndpoint(
	ctx context.Context, 
	mux *runtime.ServeMux, 
	endpoint string, 
	opts []grpc.DialOption,
) error {
	return pb.RegisterMyServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}

// Реализуем другие методы интерфейса App
```

### Обработка ошибок

Используйте обработчик ошибок для единообразных ответов:

```go
import "github.com/cripplemymind9/go-utils/server"

// При создании Runner
runner := gokit.NewRunner(
	WithErrorHandler(server.ErrorHandler()),
)
```

### Конфигурация

Сервис использует следующие переменные окружения:
- `HTTP_PORT` - порт для HTTP сервера (по умолчанию: 8080)
- `GRPC_PORT` - порт для gRPC сервера (по умолчанию: 9090)

Если переменные окружения не установлены, сервер запустится на портах по умолчанию.

```go
// Пример настройки кастомной конфигурации
config := gokit.Config{
	HTTPPort: 8080,
	GRPCPort: 9090,
}

// Применяем кастомную конфигурацию
runner := gokit.NewRunner(
	gokit.WithConfig(config),
)
```

Для запуска с кастомными портами задайте переменные окружения::
```bash
export HTTP_PORT=3000
export GRPC_PORT=3001
./your-service
```