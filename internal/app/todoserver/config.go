package todoserver

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config ...
type Config struct {
	BindAddr string `yaml:"bind_addr"`
	LogLevel string `yaml:"log_level"`
}

// NewConfig ...
func NewConfig() *Config {
	return &Config{
		BindAddr: ":8080",
		LogLevel: "debug",
	}
}

// Записывает параметры конфигурации сервера из yaml-файла
// в *Config struct при настройке сервера
func (c *Config) GetConf(confPath string) error {
	yamlFile, err := os.ReadFile(confPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return err
	}

	return err
}
