package store

import "github.com/grafyu/todo-app/model"

// UserRepository ...
type TaskRepository interface {
	Create(*model.Task) error
	FindByDate(string) (*model.Task, error)
}
