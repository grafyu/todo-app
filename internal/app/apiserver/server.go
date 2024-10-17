package apiserver

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/grafyu/todo-app/internal/app/store"
	"github.com/grafyu/todo-app/model"
)

type server struct {
	router *http.ServeMux
	logger *slog.Logger
	store  store.Store
}

// newServer - ...
func newServer(store store.Store, logger *slog.Logger) *server {
	s := &server{
		router: http.NewServeMux(),
		logger: logger,
		store:  store,
	}

	s.configureRouter()

	return s
}

// ServeHTTP() - ...
func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

// configureRouter - ...
func (s *server) configureRouter() {

	s.router.Handle("/", http.FileServer(http.Dir("./web")))

	s.router.HandleFunc("/api/nextdate", s.handleNextDate())

	s.router.HandleFunc("/api/task", s.handleTask())

	s.router.HandleFunc("/api/tasks", s.handleTaskList())

	s.router.HandleFunc("/api/task/done", s.handleCheck())

}

func (s *server) handleTask() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var (
			buf      bytes.Buffer
			resp     []byte
			answer   any
			task     model.Task
			urlValue url.Values
			err      error
		)

		taskRepo := s.store.Task()

		if r.Method == http.MethodPost || r.Method == http.MethodPut {
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			}
		} else {
			urlValue = r.URL.Query()
		}

		// Создание, Поиск, Редактирование, Удаление на пути "/api/task"
		switch r.Method {

		case http.MethodPost:
			id, err := AddTask(taskRepo, buf)
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			} else {
				answer = map[string]string{"id": id}
			}

		case http.MethodGet:
			task, err = FindTask(taskRepo, urlValue)
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			} else {
				answer = task
			}

		case http.MethodPut:
			_, err = EditTask(taskRepo, buf)
			if err != nil {
				answer = map[string]string{"error": err.Error()}

			} else {
				answer = model.Task{}
			}

		case http.MethodDelete:
			err := DeleteTask(taskRepo, urlValue)
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			} else {
				answer = map[string]string{}
			}

		}

		resp, err = json.Marshal(answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// Новая задача
func AddTask(taskRepo store.TaskRepository, buf bytes.Buffer) (string, error) {
	var task model.Task

	if err := json.Unmarshal(buf.Bytes(), &task); err != nil {
		return "", err
	}

	err := taskRepo.Create(&task)
	if err != nil {
		return "", err
	}

	return strconv.Itoa(task.ID), nil

}

func EditTask(taskRepo store.TaskRepository, buf bytes.Buffer) (model.Task, error) {

	type Temp struct {
		Id      string
		Date    string
		Title   string
		Comment string
		Repeat  string
	}

	var (
		tempo Temp
		task  model.Task
		err   error
	)

	if err = json.Unmarshal(buf.Bytes(), &tempo); err != nil {
		return model.Task{}, err
	}

	task.ID, err = strconv.Atoi(tempo.Id)
	if err != nil {
		return model.Task{}, err
	}

	task.Date = tempo.Date
	task.Title = tempo.Title
	task.Comment = tempo.Comment
	task.Repeat = tempo.Repeat

	err = taskRepo.ChangeTask(task)
	if err != nil {
		return model.Task{}, err
	} else {
		_, err := taskRepo.FindByID(task.ID)
		if err != nil {
			return model.Task{}, err
		}
	}

	return model.Task{}, nil
}

func FindTask(taskRepo store.TaskRepository, urlValue url.Values) (model.Task, error) {
	id, err := strconv.Atoi(urlValue.Get("id"))

	if err != nil {
		return model.Task{}, err
	}

	task, err := taskRepo.FindByID(id)
	if err != nil {
		return model.Task{}, err
	}

	return task, nil
}

func DeleteTask(taskRepo store.TaskRepository, urlValue url.Values) error {
	id := urlValue.Get("id")

	err := taskRepo.DeleteByID(id)
	if err != nil {
		return err
	}
	return nil
}

// вывод на экран всех задач
func (s *server) handleTaskList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var (
			answer any
			tasks  []model.Task
			err    error
		)
		taskRepo := s.store.Task()

		if r.Method == http.MethodGet {
			tasks, err = taskRepo.View()
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			} else if tasks == nil {
				answer = map[string]string{}
			} else {
				answer = map[string][]model.Task{"tasks": tasks}
			}

		}

		resp, err := json.Marshal(answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// следующая дата Периодической Задачи
func (s *server) handleNextDate() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var answer string

		now := r.URL.Query().Get("now")
		date := r.URL.Query().Get("date")
		repeat := r.URL.Query().Get("repeat")

		nowTime, err := time.Parse("20060102", now)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		answer, err = model.NextDate(nowTime, date, repeat)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		io.WriteString(w, answer)
	}
}

// Завершение выполненной задачи
func (s *server) handleCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var buf bytes.Buffer
		var (
			task   model.Task
			resp   []byte
			answer any
		)

		taskRepo := s.store.Task()

		if r.Method == http.MethodPost {
			id := r.URL.Query().Get("id")

			idStr, err := strconv.Atoi(id)
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			}

			// Ищем задачу по Id
			task, err = taskRepo.FindByID(idStr)
			if err != nil {
				answer = map[string]string{"error": err.Error()}
			}

			// удаляем если не заданы повторения
			if task.Repeat == "" {
				err := taskRepo.DeleteByID(id)
				if err != nil {
					answer = map[string]string{"error": err.Error()}
				} else {
					answer = map[string]string{}
				}

			} else {
				// ищем следующую дату задания, если повтрения заданы
				task.Date, err = model.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					answer = map[string]string{"error": err.Error()}
				}

				// устанавливаем новую дату задачи
				err = taskRepo.ChangeTask(task)
				if err != nil {
					answer = map[string]string{"err": err.Error()}
				} else {
					answer = map[string]string{}
				}
			}

		}

		resp, err := json.Marshal(answer)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}

// answer = map[string]string{"error": err.Error()}
// } else {
// 	answer = map[string]string{"id": id}
