package store

// Config - структура для хранения конфирурации store
type Config struct {
	DatabaseURL string `yaml:"database_url"`
}

// NewConfig -function helper for creating config *Config
func NewConfig() *Config {
	return &Config{}
}
