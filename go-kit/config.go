package gokit

import "github.com/spf13/viper"

const (
	httpPortKey string = "HTTP_PORT"
	grpcPortKey string = "GRPC_PORT"
)

type Config struct {
	HTTPPort int
	GRPCPort int
}

// GetConfig возвращает конфигурацию.
// Если значения не установлены через переменные окружения HTTP_PORT и GRPC_PORT,
// используются значения по умолчанию: HTTP_PORT=8080 и GRPC_PORT=9090.
func GetConfig(v *viper.Viper) Config {
	v.AutomaticEnv()

	cfg := Config{
		HTTPPort: 8080,
		GRPCPort: 9090,
	}

	if v.IsSet(httpPortKey) {
		cfg.HTTPPort = v.GetInt(httpPortKey)
	}

	if v.IsSet(grpcPortKey) {
		cfg.GRPCPort = v.GetInt(grpcPortKey)
	}

	return cfg
}
