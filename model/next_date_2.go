package model

import (
	"strings"
	"time"
)

const (
	dayInWeek   = 7
	monthInYear = 12
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	var (
		start time.Time // start task time
		next  time.Time // next task time (after current)
		err   error
	)

	loc, _ := time.LoadLocation("Europe/Moscow")

	start, err = time.ParseInLocation("20060102", date, loc)
	if err != nil {
		return "", err
	}

	now = now.Round(1h0m0s)

	switch strings.Split(repeat)[1] {
	case "y":

	}

}
