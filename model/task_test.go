package model_test

import (
	"testing"

	"github.com/grafyu/todo-app/model"
	"github.com/stretchr/testify/assert"
)

func TestTask_VAlidate(t *testing.T) {
	tsk := model.TestTask(t)
	assert.NoError(t, tsk.Validate())
}

func TestTask_Beforcreate(t *testing.T) {
	tsk := model.TestTask(t)
	assert.NoError(t, tsk.BeforeCreate())
}
