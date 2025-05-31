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

Настройка портов и других параметров:

```go
// Переопределяем конфигурацию по умолчанию
config := gokit.Config{
	HTTPPort: 8080,
	GRPCPort: 9090,
}

// Применяем кастомную конфигурацию
runner := gokit.NewRunner(
	gokit.WithConfig(config),
)
```