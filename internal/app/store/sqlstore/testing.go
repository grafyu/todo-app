package sqlstore

import (
	"database/sql"
	"fmt"
	"testing"
)

// TestDB - возвращает сконфигурированный тестовый Store,
// а, также функцию, которая будет очищать заполненные
// в процессе тестов для проведения последующих тестов
func TestDB(t *testing.T, databaseURL string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("sqlite", databaseURL)
	if err != nil {
		t.Fatal(err)
	}

	if err := CreateTable(db, databaseURL); err != nil {
		t.Fatal(err)
	}

	return db, func(tables ...string) {
		if len(tables) > 0 {
			for _, table := range tables {
				if _, err := db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
					t.Fatal(err)
				}

			}
		}
		db.Close()
	}
}
