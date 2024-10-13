package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Task - представление записей Task из DB
type Task struct {
	ID      int    `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

func checkRepeatFormat(repeat interface{}) error {
	rule := strings.Split(repeat.(string), " ")

	// errFormat := errors.New("not correct the repeate format")

	// symbol check
	if !slices.Contains([]string{"y", "d", "w", "m", ""}, rule[0]) {
		return errors.New("invalid symbol")
	}

	// rule string format check
	switch rule[0] {
	case "y":
		if len(rule) != 1 {
			return errors.New("not correct the repeate year format")
		}

		return nil

	case "d":
		if len(rule) < 2 {
			return errors.New("not specified interval in days")
		} else if len(rule) > 2 {
			return errors.New("not correct the repeate day format")
		}

		interval, err := strconv.Atoi(rule[1])
		if err != nil {
			return err
		}

		if interval > 400 || interval < 1 {
			return errors.New("wrong number of days in an interval")
		}

		return nil

	case "w":
		// check items
		if len(rule) < 2 {
			return errors.New("not specified days of the week")
		} else if len(rule) > 2 {
			return errors.New("not correct the repeate week format")
		}

		// check rules param
		daysWeek := strings.Split(rule[1], ",")
		numsWeek, err := charsToInts(daysWeek)
		if err != nil {
			return err
		}

		slices.Sort(numsWeek)
		if numsWeek[0] < 1 || numsWeek[len(numsWeek)-1] > 7 {
			return errors.New("wrong number of days of the week")
		}

		return nil

	case "m":
		if len(rule) < 2 {
			return errors.New("not specified days of month")
		}

		// check rules param
		if len(rule) == 2 {
			daysMnth := strings.Split(rule[1], ",")
			numsMnth, err := charsToInts(daysMnth)
			if err != nil {
				return err
			}

			slices.Sort(numsMnth)
			if numsMnth[0] < -2 || numsMnth[len(numsMnth)-1] > 31 {
				return errors.New("wrong number days of the month")
			}
		}

		if len(rule) == 3 {
			daysMnth := strings.Split(rule[2], ",")
			numsMnth, err := charsToInts(daysMnth)
			if err != nil {
				return err
			}

			slices.Sort(numsMnth)
			if numsMnth[0] < 1 || numsMnth[len(numsMnth)-1] > 12 {
				return errors.New("wrong number month")
			}
		}

	}
	return nil
}

func checkDateFormat(date interface{}) error {
	dateField := date.(string)

	if dateField == "" {
		return errors.New("не задана дата задания")
	}

	dateTime, err := time.Parse("20060102", dateField)
	if err != nil {
		return errors.New("неверный формат даты задания")
	}

	if dateTime.Before(time.Now()) {
		return errors.New("дата задания не может быть раньше текущей даты")
	}

	return nil

}

// Validate() - валидация данных task
// отправленных frontend
func (tsk *Task) Validate() error {
	return validation.ValidateStruct(
		tsk,
		validation.Field(&tsk.Title, validation.Required),
		validation.Field(&tsk.Date, validation.By(checkDateFormat)),
		// validation.Field(&tsk.Date, validation.Date("20060102")),
		validation.Field(&tsk.Repeat, validation.By(checkRepeatFormat)),
	)

}

// BeforeCreate - запускается при созданием task
// здесь валидируются данные перед запись task в DB
func (tsk *Task) BeforeCreate() error {
	return nil
}

// После возвращения ответа на запрос от handler
func BeforeSend() error {
	return nil
}

func (t Task) MarshalJSON() ([]byte, error) {

	var id string = ""

	id = fmt.Sprintf("%d", t.ID)
	aux := struct {
		ID      string `json:"id"`
		Date    string `json:"date"`
		Title   string `json:"title"`
		Comment string `json:"comment"`
		Repeat  string `json:"repeat"`
	}{
		ID:      id,
		Date:    t.Date,
		Title:   t.Title,
		Comment: t.Comment,
		Repeat:  t.Repeat,
	}

	return json.Marshal(aux)
}

// func (t Task) UnmarshalJSON([]byte) error {

// 	var id int

// 	t.ID, err := strconv.Atoi(id)
// 	if err != nil {
// 		return err
// 	}

// 	aux := struct {
// 		ID      int    `json:"id"`
// 		Date    string `json:"date"`
// 		Title   string `json:"title"`
// 		Comment string `json:"comment"`
// 		Repeat  string `json:"repeat"`
// 	}{
// 		ID:      id,
// 		Date:    t.Date,
// 		Title:   t.Title,
// 		Comment: t.Comment,
// 		Repeat:  t.Repeat,
// 	}

// 	return nil
// }
