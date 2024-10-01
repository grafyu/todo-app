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
	err := config.GetConf(configPath)
	if err != nil {
		log.Fatal(err)
	}

	s := apiserver.New(config)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
