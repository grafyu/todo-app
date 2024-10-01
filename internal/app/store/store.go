package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	_ "modernc.org/sqlite"
)

// Store - инкапсулирует детали реализации взаимодействия
// с DB при помощи public method
// ... repository
type Store struct {
	config         *Config
	db             *sql.DB
	taskRepository *TaskRepository
}

// New() creates the object “store” according to the
// given configuration “config”
func New(config *Config) *Store {
	// func New(db *sql.DB) *Store {
	// ...
	return &Store{
		config: config,
		// db:     db,
	}
}

// CreateTable() - create table if fileDB exist
func CreateTable(db *sql.DB, dbURL string) error {
	queries := []string{`CREATE TABLE IF NOT EXISTS scheduler (
		id integer primary key autoincrement,
		date char(8) not null default "",
		title varchar(128) not null default "",
		comment text not null default "",
		repeat varchar(128) not null default "")`,
		`CREATE INDEX scheduler_date ON scheduler (date)`}

	appPath, err := os.Executable()
	if err != nil {
		log.Printf("Error %s when finding resource located relative to an executable", err)
		return err
	}

	dbFile := filepath.Join(filepath.Dir(appPath), dbURL)
	_, err = os.Stat(dbFile)

	if err != nil {
		fmt.Println("Creating a new database")
		ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancelfunc()

		for _, query := range queries {
			_, err := db.ExecContext(ctx, query)
			if err != nil {
				log.Printf("Error %s when creating 'scheduler' table", err)
				return err
			}
		}
	}
	return nil
}

// Open() - open DB
func (s *Store) Open() error {
	db, err := sql.Open("sqlite", s.config.DatabaseURL)
	if err != nil {
		return err
	}

	if err := CreateTable(db, s.config.DatabaseURL); err != nil {
		return err
	}

	s.db = db

	return nil
}

// Close() - disconnects the application from DB with necessary actions
func (s *Store) Close() {
	s.db.Close()
}

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
