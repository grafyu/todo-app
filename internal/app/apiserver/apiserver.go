package apiserver

import (
	"database/sql"
	"log"
	"log/slog"
	"net/http"
	"os"

	"github.com/grafyu/todo-app/internal/app/store/sqlstore"
)

func Start(config *Config) error {
	db, err := newDB(config.DatabaseURL)
	if err != nil {
		return err
	}

	defer db.Close()
	store := sqlstore.New(db)
	logger := newLogger(config.LogLevel)
	srv := newServer(store, logger)

	return http.ListenAndServe(config.BindAddr, srv)
}

func newDB(databaseURL string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		return nil, err
	}

	if err := sqlstore.CreateTable(db, databaseURL); err != nil {
		return nil, err
	}

	return db, nil
}

func newLogger(logLevel string) *slog.Logger {
	levelVar := new(slog.Level)

	if err := levelVar.UnmarshalText([]byte(logLevel)); err != nil {
		log.Println(err.Error())
		levelVar = nil
	}

	textHandler := slog.NewTextHandler(
		os.Stderr,
		&slog.HandlerOptions{Level: levelVar})

	return slog.New(textHandler)
}
