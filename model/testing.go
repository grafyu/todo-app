// готовая модель для проведения тестов

package model

import "testing"

func TestTask(tu *testing.T) *Task {
	return &Task{
		ID:      1000,
		Date:    "20201231",
		Title:   "Test task",
		Comment: "Task for test",
		Repeat:  "d 5",
	}
}
