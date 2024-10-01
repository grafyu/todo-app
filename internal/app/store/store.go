package store

import (
	"database/sql"

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

// // Open() - open DB
// func (s *Store) Open() error {
// 	db, err := sql.Open("sqlite", s.config.DatabaseURL)
// 	if err != nil {
// 		return err
// 	}

// 	if err := CreateTable(db, s.config.DatabaseURL); err != nil {
// 		return err
// 	}

// 	s.db = db

// 	return nil
// }

// // Close() - disconnects the application from DB with necessary actions
// func (s *Store) Close() {
// 	s.db.Close()
// }

// Task() - метод объкта типа Store для получения объекта TaskRepository
// для работы с ним
func (s *Store) Task() *TaskRepository {
	if s.taskRepository != nil {
		return s.taskRepository
	}

	s.taskRepository = &TaskRepository{
		store: s,
	}

	return s.taskRepository
}
