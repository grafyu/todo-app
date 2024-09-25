package todoserver

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var lvlVar = new(slog.LevelVar)

type ToDoServer struct {
	config *Config
	logger *slog.Logger
	router *http.ServeMux
}

// New ToDoServer ...
func New(config *Config) *ToDoServer {
	return &ToDoServer{
		config: config,
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvlVar})),
		router: http.NewServeMux(),
	}
}

// Start HTTP server, connect to DB and etc.
func (s *ToDoServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}
	// s.logger.Debug("no error TO-DO server")
	s.configureRouter()

	return http.ListenAndServe(s.config.BindAddr, s.router) // return err = http.ListenAndServe()
}

// configureLogger - configures the logger
func (s *ToDoServer) configureLogger() error {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(s.config.LogLevel))
	if err != nil {
		return err
	}

	fmt.Printf("log level is %v\n", level)

	lvlVar.Set(*level) // set the level from todoserver.yaml
	return nil
}

func (s *ToDoServer) configureRouter() {
	s.router.Handle("/", http.FileServer(http.Dir("./web")))
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *ToDoServer) handleHello() http.HandlerFunc {
	// ...

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
