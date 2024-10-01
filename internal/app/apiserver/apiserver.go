package apiserver

import (
	"database/sql"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/grafyu/todo-app/internal/app/store/sqlstore"
)

var lvlVar = new(slog.LevelVar)

type APIServer struct {
	config *Config
	logger *slog.Logger
	router *http.ServeMux
	store  *sqlstore.Store
}

// New APIServer ...
func New(config *Config) *APIServer {
	return &APIServer{
		config: config,
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvlVar})),
		router: http.NewServeMux(),
	}
}

// Start HTTP server, connect to DB and etc.
func (s *APIServer) Start() error {
	if err := s.configureLogger(); err != nil {
		return err
	}

	s.configureRouter()

	// конфигурирует Store, вызывает у него метод Open()
	// если все удается, записать в поле `store` ToDoserver ссылку на наше хранилище
	if err := s.configureStore(); err != nil {
		return err
	}

	return http.ListenAndServe(s.config.BindAddr, s.router) // return err = http.ListenAndServe()
}

// configureLogger - configures the logger
func (s *APIServer) configureLogger() error {
	level := new(slog.Level)
	err := level.UnmarshalText([]byte(s.config.LogLevel))
	if err != nil {
		return err
	}

	fmt.Printf("log level is %v\n", level)

	lvlVar.Set(*level) // set the level from todoserver.yaml
	return nil
}

func (s *APIServer) configureRouter() {
	s.router.Handle("/", http.FileServer(http.Dir("./web")))
	s.router.HandleFunc("/hello", s.handleHello())
}

func (s *APIServer) configureStore() error {
	db, err := sql.Open("sqlite", s.config.Store)
	if err != nil {
		return err
	}

	if err := sqlstore.CreateTable(db, s.config.Store); err != nil {
		return err
	}

	s.store = sqlstore.New(db)

	return nil
}

func (s *APIServer) handleHello() http.HandlerFunc {
	// ...

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
