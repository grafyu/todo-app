package store

// Store ...
type Store interface {
	Task() TaskRepository
}
