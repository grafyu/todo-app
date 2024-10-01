package model

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	validation "github.com/go-ozzo/ozzo-validation"
)

// Task - представление записей Task из DB
type Task struct {
	ID      int
	Date    string
	Title   string
	Comment string
	Repeat  string
}

// checkParam() - проверяем значения параметров
// по границам заданного диапазона
func checkParam(values string, min int, max int) error {
	for _, value := range strings.Split(values, ",") { // разбиваем группу аргументов на отдельные значения
		num, _ := strconv.Atoi(value)
		if num < min || num > max || num == 0 {
			return errors.New("недопустимое значение параметров повторения\nвыход за диапазон или нулевое значение")
		}
	}

	return nil
}

func checkRepeatFormat(repeat interface{}) error {
	// в зависимости от первого символа здесь будет
	// выбираться шаблон для валидации
	params := strings.Split(repeat.(string), " ")

	switch params[0] {
	case "y":
		if len(params) == 1 {
			return nil
		}
	case "d":
		switch len(params) {
		case 1:
			return errors.New("не указан интервал в днях")
		case 2:
			return checkParam(params[1], 1, 400)
		default:
			return errors.New("too many parameters")
		}

	case "w":
		switch len(params) {
		case 1:
			return errors.New("не указаны дни недели")
		case 2:
			return checkParam(params[1], 1, 12)
		default:
			return errors.New("too many parameters")
		}
	case "m":
		switch len(params) {
		case 1:
			return errors.New("не указаны дни месяца")
		case 2:
			return checkParam(params[1], -2, 31)
		case 3:
			if err := checkParam(params[1], -2, 31); err != nil {
				return err
			}
			return checkParam(params[2], 1, 12)
		default:
			return errors.New("too many parameters")
		}
	default:
		return errors.New("недопустимый символ параметра")

	}

	fmt.Println("All right !!!")

	return nil
}

// Validate() - валидация данных task
// отправленных frontend
func (tsk *Task) Validate() error {
	return validation.ValidateStruct(
		tsk,
		validation.Field(&tsk.Title, validation.Required),
		validation.Field(&tsk.Date, validation.Required, validation.Date("20060102")),
		validation.Field(&tsk.Repeat, validation.By(checkRepeatFormat)),
	)

}

// BeforeCreate - запускается при созданием task
// здесь валидируются данные перед запись task в DB
func (tsk *Task) BeforeCreate() error {
	// if len(t.Title) > 0
	return nil
}

// После возвращения ответа на запрос от handler
func BeforeSend() error {
	return nil
}

// NextDate - calculates the next date according to the specified rule
func NextDate(now time.Time, date string, repeat string) (string, error) {
	return "", nil
}
