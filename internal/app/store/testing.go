package store

import (
	"fmt"
	"testing"
)

// TestStore - возвращает сконфигурированный тестовый Store,
// а, также функцию, которая будет очищать заполненные
// в процессе тестов для проведения последующих тестов
func TestStore(t *testing.T, databaseURL string) (*Store, func(...string)) {
	t.Helper()

	config := NewConfig()
	config.DatabaseURL = databaseURL
	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s, func(tables ...string) {
		if len(tables) > 0 {
			for _, table := range tables {
				if _, err := s.db.Exec(fmt.Sprintf("DROP TABLE IF EXISTS %s", table)); err != nil {
					t.Fatal(err)
				}

			}
		}
		s.Close()
	}
}
