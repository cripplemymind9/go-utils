package gokit

// Config содержит конфигурацию для Runner
type Config struct {
	HTTPPort int
	GRPCPort int
}

// GetConfig возвращает конфигурацию по умолчанию
func GetConfig() Config {
	return Config{
		HTTPPort: 8080,
		GRPCPort: 9090,
	}
}
