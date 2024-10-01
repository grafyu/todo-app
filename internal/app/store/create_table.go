package store

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"
)

// CreateTable() - create Table if FileDB NOT_EXIST
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
