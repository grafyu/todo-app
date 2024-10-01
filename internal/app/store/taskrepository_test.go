package store_test

import (
	"testing"

	"github.com/grafyu/todo-app/internal/app/store"
	"github.com/grafyu/todo-app/model"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := store.TestDB(t, databaseURL)
	defer teardown("scheduler")

	s := store.New(db)
	tsk, err := s.Task().Create(model.TestTask(t))
	assert.NoError(t, err)
	assert.NotNil(t, tsk)
}
 