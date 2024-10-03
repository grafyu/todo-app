package main

import (
	"flag"
	"log"

	"github.com/grafyu/todo-app/internal/app/apiserver"
)

var (
	configPath string
)

func init() {
	flag.StringVar(&configPath, "config-path", "configs/todoserver.yaml", "path to config file")
}

func main() {
	flag.Parse()

	config := apiserver.NewConfig()
	err := config.YamlToConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	if err := apiserver.Start(config); err != nil {
		log.Fatal(err)
	}
}
