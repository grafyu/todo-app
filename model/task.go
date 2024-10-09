package model

import (
	"errors"
	"slices"
	"strconv"
	"strings"

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

func checkRepeatFormat(repeat interface{}) error {
	params := strings.Split(repeat.(string), " ")

	errFormat := errors.New("not correct the repeate format")

	// symbol check
	if !slices.Contains([]string{"y", "d", "w", "m"}, params[0]) {
		return errors.New("invalid symbol")
	}

	// rule string format check
	switch params[0] {
	case "y":
		if len(params) > 1 {
			return errFormat
		}
	case "d":
		if len(params) < 2 {
			return errors.New("missing interval in days")
		}

		if len(params) == 2 {
			days, err := strconv.Atoi(params[1])
			if err != nil {
				return err
			}

			if days < 1 && days > 400 {
				return errors.New("maximum allowable interval is exceeded")
			}
		}

		return errFormat

	case "w":
		if len(params) < 2 {
			return errors.New("missing days of the week")
		}

		if len(params) == 2 {
			weekdays, err := charToInt(strings.Split(params[1], ","))
			if err != nil {
				return err
			}

			if slices.Min(weekdays) < 1 || slices.Max(weekdays) > 7 {
				return errors.New("invalid value of weekday")
			}
		}

		return errFormat

	case "m":
		if len(params) < 2 {
			return errors.New("missing days of the month")
		}

		if len(params) == 2 {
			monthdays, err := charToInt(strings.Split(params[1], ","))
			if err != nil {
				return err
			}

			if slices.Min(monthdays) < -2 || slices.Max(monthdays) > 31 || slices.Max(monthdays) == 0 {
				return errors.New("invalid value of monthday")
			}
		}

		if len(params) == 3 {
			months, err := charToInt(strings.Split(params[2], ","))
			if err != nil {
				return err
			}

			if slices.Min(months) < 1 || slices.Max(months) > 12 {
				return errors.New("invalid value of month")
			}
		}

		return errFormat
	}

	return nil

}

// convert []string to []int and sort result slice
func charToInt(charSl []string) ([]int, error) {
	var intSl []int
	for _, char := range charSl {
		num, err := strconv.Atoi(char)
		if err != nil {
			return nil, err
		}

		intSl = append(intSl, num)
	}

	slices.Sort(intSl)

	return intSl, nil
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
