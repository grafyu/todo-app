// репозитории - отвечают за взаимодействие моделей с DB

package store

import (
	"database/sql"

	"github.com/grafyu/todo-app/model"
)

// UserRepository
type TaskRepository struct {
	store *Store
}

// Create ...
func (r TaskRepository) Create(t *model.Task) (*model.Task, error) {
	if err := t.Validate(); err != nil {
		return nil, err
	}

	if err := t.BeforeCreate(); err != nil {
		return nil, err
	}

	if err := r.store.db.QueryRow(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat) RETURNING id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	).Scan(&t.ID); err != nil {
		return nil, err
	}
	return t, nil
}

// FindByDate - найти Task по дате,
func (r TaskRepository) FindByDate(date string) (*model.Task, error) {
	return nil, nil
}
