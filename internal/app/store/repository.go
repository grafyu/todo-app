package store

import "github.com/grafyu/todo-app/model"

// UserRepository ...
type TaskRepository interface {
	Create(*model.Task) error
	View() ([]model.Task, error)
	FindByDate(string) (*model.Task, error)
	FindByID(int) (*model.Task, error)
	ChangeTask(*model.Task) error
	DeleteByID(string) error
	// FindByRule(string) (*model.Task, error)
}
