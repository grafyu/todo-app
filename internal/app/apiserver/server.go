package apiserver

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/grafyu/todo-app/internal/app/store"
)

type server struct {
	router *http.ServeMux
	logger *slog.Logger
	store  store.Store
}

// newServer - ...
func newServer(store store.Store) *server {
	s := &server{
		router: http.NewServeMux(),
		logger: slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: lvlVar})),
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
}
