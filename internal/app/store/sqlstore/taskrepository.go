// репозитории - отвечают за взаимодействие моделей с DB

package sqlstore

import (
	"database/sql"

	"github.com/grafyu/todo-app/internal/app/store"
	"github.com/grafyu/todo-app/model"
)

// UserRepository
type TaskRepository struct {
	store *Store
}

// Create ...
func (r TaskRepository) Create(t *model.Task) error {
	if err := t.Validate(); err != nil {
		return err
	}

	if err := t.BeforeCreate(); err != nil {
		return err
	}

	return r.store.db.QueryRow(
		"INSERT INTO scheduler (date, title, comment, repeat) VALUES (:date, :title, :comment, :repeat) RETURNING id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
	).Scan(&t.ID)
}

// FindByDate - найти Task по дате,
func (r TaskRepository) FindByDate(date string) (*model.Task, error) {
	tsk := &model.Task{}
	if err := r.store.db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE date = :date",
		sql.Named("date", date),
	).Scan(
		&tsk.ID,
		&tsk.Date,
		&tsk.Title,
		&tsk.Comment,
		&tsk.Repeat,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecondNotFound
		}
		
		return nil, err
	}

	return tsk, nil
}
