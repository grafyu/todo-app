package apiserver

import (
	"log/slog"
	"net/http"

	"github.com/grafyu/todo-app/internal/app/store"
)

type server struct {
	router *http.ServeMux
	logger *slog.Logger
	store  store.Store
}

// newServer - ...
func newServer(store store.Store, logger *slog.Logger) *server {
	s := &server{
		router: http.NewServeMux(),
		logger: logger,
		store:  store,
	}

	s.configureRouter()
	return s
}

// ServeHTTP() - ...
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter - ...
func (s *server) configureRouter() {
	// ...
	s.router.Handle("/", http.FileServer(http.Dir("./web")))
	// s.router.HandleFunc("")
}
