package todoserver

import (
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"

	"github.com/grafyu/todo-app/internal/app/store"
)

var lvlVar = new(slog.LevelVar)

type ToDoServer struct {
	config *Config
	logger *slog.Logger
	router *http.ServeMux
	store  *store.Store
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

	// конфигурирует Store, вызывает у него метод Open()
	// если все удается, записать в поле `store` ToDoserver ссылку на наше хранилище
	s.configureRouter()

	if err := s.configureStore(); err != nil {
		return err
	}

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

func (s *ToDoServer) configureStore() error {
	st := store.New(s.config.Store)   // создаем новое хранилище, по config-ам
	if err := st.Open(); err != nil { // открываем новое хранилище
		return err
	}

	s.store = st

	return nil
}

func (s *ToDoServer) handleHello() http.HandlerFunc {
	// ...

	return func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello")
	}
}
