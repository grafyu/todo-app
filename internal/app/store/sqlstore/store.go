package sqlstore

import (
	"database/sql"

	"github.com/grafyu/todo-app/internal/app/store"
	_ "modernc.org/sqlite"
)

// Store - инкапсулирует детали реализации взаимодействия
// с DB при помощи public method
// ... repository
type Store struct {
	db             *sql.DB
	taskRepository *TaskRepository
}

// New() - creates the object “store”
func New(db *sql.DB) *Store {
	// ...
	return &Store{
		db: db,
	}
}

// Task() - метод объкта типа Store для получения объекта TaskRepository
// для работы с ним
func (s *Store) Task() store.TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
	}

	return s.taskRepository
}
