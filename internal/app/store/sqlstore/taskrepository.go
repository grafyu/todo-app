// репозитории - отвечают за взаимодействие моделей с DB

package sqlstore

import (
	"database/sql"
	"strconv"

	"github.com/grafyu/todo-app/internal/app/store"
	"github.com/grafyu/todo-app/model"
)

// TaskRepository ...
type TaskRepository struct {
	store *Store
}

// Create Task...
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

// View Tasks...
func (r TaskRepository) View() ([]model.Task, error) {
	var (
		tasks []model.Task
		tsk   model.Task
	)

	rows, err := r.store.db.Query("SELECT id, date, title, comment, repeat FROM scheduler ORDER BY date LIMIT 25")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&tsk.ID,
			&tsk.Date,
			&tsk.Title,
			&tsk.Comment,
			&tsk.Repeat,
		)
		if err != nil {
			return nil, err
		}

		tasks = append(tasks, tsk)

		if err := rows.Err(); err != nil {
			return nil, err
		}

	}

	return tasks, nil
}

// FindByDate - найти Task по дате,
func (r TaskRepository) FindByDate(date string) (model.Task, error) {
	tsk := model.Task{}

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
			return model.Task{}, store.ErrRecondNotFound
		}
	}

	return tsk, nil
}

func (r TaskRepository) ChangeTask(t model.Task) error {

	if err := t.Validate(); err != nil {
		return err
	}

	if err := t.BeforeCreate(); err != nil {
		return err
	}

	_, err := r.store.db.Exec("UPDATE scheduler SET date = :date, title = :title, comment = :comment, repeat = :repeat WHERE id = :id",
		sql.Named("date", t.Date),
		sql.Named("title", t.Title),
		sql.Named("comment", t.Comment),
		sql.Named("repeat", t.Repeat),
		sql.Named("id", t.ID),
	)

	return err
}

func (r TaskRepository) DeleteByID(id string) error {
	idStr, err := strconv.Atoi(id)
	if err != nil {
		return err
	}

	r.store.db.QueryRow("DELETE FROM scheduler WHERE id = :id ",
		sql.Named("id", idStr),
	)

	return nil
}

// FindByDate - найти Task по дате,
func (r TaskRepository) FindByID(id int) (model.Task, error) {
	tsk := model.Task{}

	if err := r.store.db.QueryRow(
		"SELECT id, date, title, comment, repeat FROM scheduler WHERE id = :id",
		sql.Named("id", id),
	).Scan(
		&tsk.ID,
		&tsk.Date,
		&tsk.Title,
		&tsk.Comment,
		&tsk.Repeat,
	); err != nil {
		if err == sql.ErrNoRows {
			return model.Task{}, err
		}
	}

	return tsk, nil
}
