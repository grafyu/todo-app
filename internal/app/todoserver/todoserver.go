package todoserver

import (
	"fmt"
	"log/slog"
	"os"
)

var lvlVar = new(slog.LevelVar)

type ToDoServer struct {
	config *Config
	logger *slog.Logger
}

// New ToDoServer ...
func New(config *Config) *ToDoServer {
	return &ToDoServer{
		config: config,
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvlVar})),
	}
}

// Start HTTP server, connect to DB and etc.
func (s *ToDoServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.logger.Debug("no error TO-DO server")
	s.logger.Info("starting TO-DO server")

	return nil
}

// configureLogger - configures the logger
func (s *ToDoServer) configureLogger() error {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(s.config.LogLevel))
	if err != nil {
		return err
	}

	fmt.Printf("log level is %v\n\n", level)

	lvlVar.Set(*level) // set the level from todoserver.yaml
	return nil
}
