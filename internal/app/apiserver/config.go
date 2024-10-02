package apiserver

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

// Config - struct for storing the TO-DO server configuration
type Config struct {
	BindAddr    string `yaml:"bind_addr"`
	LogLevel    string `yaml:"log_level"`
	DatabaseURL string `yaml:"database_url"`
}

// NewConfig - creating configuration for run ToDo Server
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

	// If an environment variable exists, BindAddr is set to its value
	addr, valid := portValidation(os.Getenv("TODO_ADDR"))
	if valid {
		c.BindAddr = ":" + addr
		fmt.Printf("The BindAddr is set to %s from the $TODO_ADDR\n", addr)
	}

	return err
}

// Checks the presence and validity of the address in the environment variable
func portValidation(val string) (string, bool) {
	if num, err := strconv.Atoi(val); err == nil && num >= 7000 && num <= 9000 {
		return val, true
	}
	return "", false
}
