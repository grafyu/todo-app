package store_test

import (
	"testing"

	"github.com/grafyu/todo-app/internal/app/store"
	"github.com/grafyu/todo-app/model"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s, teardown := store.TestStore(t, databaseURL)
	defer teardown("scheduler")

	tsk, err := s.Task().Create(model.TestTask(t))
	assert.NoError(t, err)
	assert.NotNil(t, tsk)
}
