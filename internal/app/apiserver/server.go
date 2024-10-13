package apiserver

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
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
	// ...
	s.router.Handle("/", http.FileServer(http.Dir("./web")))
	// s.router.HandleFunc("")
	s.router.HandleFunc("/api/nextdate", s.handleNextDate())

	s.router.HandleFunc("/api/task", s.handleCreateTask())

	s.router.HandleFunc("/api/tasks", s.handleTaskList())

	s.router.HandleFunc("/api/task/done", s.handleCheck())

}

func (s *server) handleCreateTask() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var task model.Task
		var buf bytes.Buffer
		var resp []byte

		answer := make(map[string]any)

		switch r.Method {
		// создание Задачи
		case http.MethodPost:
			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			err = s.store.Task().Create(&task)
			if err != nil {
				// http.Error(w, err.Error(), http.StatusInternalServerError)
				answer["error"] = err.Error()
			} else {
				answer["id"] = strconv.Itoa(task.ID)
			}

			resp, err = json.Marshal(answer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		// поиск Задачи по id
		case http.MethodGet:
			id, err := strconv.Atoi(r.URL.Query().Get("id"))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			task, err := s.store.Task().FindByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			resp, err = json.Marshal(task)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		// Редактирование Задачи
		case http.MethodPut:

			type Temp struct {
				ID      string
				Date    string
				Title   string
				Comment string
				Repeat  string
			}

			var (
				temp Temp
				task model.Task
			)

			answer := make(map[string]any)

			_, err := buf.ReadFrom(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err = json.Unmarshal(buf.Bytes(), &temp); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			task.ID, err = strconv.Atoi(temp.ID)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			task.Date = temp.Date
			task.Title = temp.Title
			task.Comment = temp.Comment
			task.Repeat = temp.Repeat

			err = s.store.Task().ChangeTask(&task)
			if err != nil {
				answer["error"] = "Задача не найдена"
			}

			resp, err = json.Marshal(answer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		// удаление Задачи
		case http.MethodDelete:
			id := r.URL.Query().Get("id")

			err := s.store.Task().DeleteByID(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			resp, err = json.Marshal(answer)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)

	}
}

// вывод на экран всех задач
func (s *server) handleTaskList() http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		answer := make(map[string]any)

		if r.Method == http.MethodGet {
			tasks, err := s.store.Task().View()
			if err != nil {
				answer["error"] = err.Error()
			} else if tasks == nil {
				answer["tasks"] = []model.Task{}
			} else {
				answer["tasks"] = tasks
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

// отметка выполненной задачи
func (s *server) handleCheck() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// var task model.Task
		// var buf bytes.Buffer
		var resp []byte

		// answer := make(map[string]any)

		if r.Method == http.MethodPost {
			id := r.URL.Query().Get("id")

			idStr, err := strconv.Atoi(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			task, err := s.store.Task().FindByID(idStr)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			if task.Repeat == "" {
				err := s.store.Task().DeleteByID(id)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

			} else {
				task.Date, err = model.NextDate(time.Now(), task.Date, task.Repeat)
				if err != nil {
					fmt.Println("here")
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				err = s.store.Task().ChangeTask(task)
				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}

		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(resp)
	}
}
