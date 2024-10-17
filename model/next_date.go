package model

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"time"
)

func NextDate(now time.Time, date string, repeat string) (string, error) {

	loc, _ := time.LoadLocation("Europe/Moscow")

	start, err := time.ParseInLocation("20060102", date, loc)
	if err != nil {
		return "", err
	}

	now = now.In(loc).Truncate(time.Hour)

	fmt.Println(now) // test

	next := start

	weekdays := map[string]string{
		"1": "Monday",
		"2": "Tuesday",
		"3": "Wednesday",
		"4": "Thursday",
		"5": "Friday",
		"6": "Saturday",
		"7": "Sunday",
	}

	months := map[string]string{
		"1":  "January",
		"2":  "February",
		"3":  "March",
		"4":  "April",
		"5":  "May",
		"6":  "June",
		"7":  "July",
		"8":  "August",
		"9":  "September",
		"10": "October",
		"11": "November",
		"12": "December",
	}

	switch strings.Split(repeat, " ")[0] {
	case "y":
		for next.Before(now) {
			next = next.AddDate(1, 0, 0)
		}

		return next.AddDate(1, 0, 0).Format("20060102"), nil

	case "d":
		interval, err := strconv.Atoi(strings.Split(repeat, " ")[1])
		if err != nil {
			return "", err
		}

		// next = start.AddDate(0, 0, interval)
		next = start

		for next.Before(now) {
			next = next.AddDate(0, 0, interval)
		}

		return next.AddDate(0, 0, interval).Format("20060102"), nil

	case "w":
		digits := strings.Split(repeat, " ")[1]
		// digits = strings.Replace(digits, "7", "0", 1)		// заменили "7" на "0"

		daysNumb := strings.Split(digits, ",")

		return nextCheckDay(next, daysNumb, weekdays).Format("20060102"), err

	case "m":
		var (
			daysNmb, monthNmb, monthPoint []string
			next                          time.Time
		)

		// if now.Before(start) {
		// 	next = start
		// } else {
		// 	next = now
		// }

		next = now

		daysNmb = relativeDay(next,
			strings.Split(strings.Split(repeat, " ")[1], ","))

		if len(strings.Split(repeat, " ")) == 3 {
			// []string
			monthNmb = strings.Split(strings.Split(repeat, " ")[2], ",")

			// составляем monthNumb - список заданных месяцев
			for _, mn := range monthNmb {
				monthPoint = append(monthPoint, months[mn])
			}

		} else {
			// либо включаем все месяца (если специально не указаны в правиле)
			for _, value := range months {
				monthPoint = append(monthPoint, value)
			}
		}

		slices.Sort(monthPoint)
		slices.Sort(daysNmb)

		next, err = findDayInMonth(next, daysNmb, monthPoint)
		if err != nil {
			return "", err
		}

		next = next.AddDate(0, 0, 1)

		next, err = findDayInMonth(next, daysNmb, monthPoint)
		if err != nil {
			return "", err
		}

		return next.Format("20060102"), nil
	}

	return "", errors.New("неправильная буква repeate rule mode")
}

// convert []string to []int and sort result slice
func charsToInts(charSl []string) ([]int, error) {
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

func nextCheckDay(next time.Time, daysNumb []string, weekdays map[string]string) time.Time {
	var days []string

	for _, dn := range daysNumb {
		days = append(days, weekdays[dn])
	}

	checkDayNow := slices.Contains(days, next.Weekday().String())

	// если данный "next" это checkDayNow, то пропускаем его поиск
	if !checkDayNow {
		for !slices.Contains(days, next.Weekday().String()) {
			next = next.AddDate(0, 0, 1)
		}
	}

	// ищем checkDayNew ...
	for !slices.Contains(days, next.Weekday().String()) {
		next = next.AddDate(0, 0, 1)
	}

	return next
}

// func nextCheckMonthday(next time.Time, monthNmb []string, daysNmb []string, month [string]string) {
// 	var months []string

// 	for mn := range monthNmb {
// 		months = append(months, weekdays[mn])
// 	}

// 	// проверить - next есть в monthNumb
// 	checkMonthNow := slices.Contains(months, int(next.Month()))

// 	if !checkMonthNow {
// 		nextMonthday()
// 	}

// 	if !checkDayNow {
// 		nextMonthday()
// 	}

// }

// ищет заданные дни в заданном месяце
func findDayInMonth(nxt time.Time, daysNmb []string, monthPoint []string) (time.Time, error) {

	// проверить, что monthPoint не пустой !!!!

	for range monthPoint {
		for !slices.Contains(monthPoint, nxt.Month().String()) {
			nxt = nxt.AddDate(0, 1, 0)
		}

		// перебираем дни до конца месяца, пока не сменится текущий месяц
		currMonth := nxt.Month()

		for currMonth == nxt.Month() {
			if slices.Contains(daysNmb, fmt.Sprint(nxt.Day())) {
				return nxt, nil
			}
			nxt = nxt.AddDate(0, 0, 1)
		}
	}

	return nxt, errors.New("не найден день по заданному rule repeate")
}


// Заменяем относительные номера дней месяца
func relativeDay(nxt time.Time, daysNmb []string) []string {

	for _, d := range []string{"-1", "-2"} {
		iDay := slices.Index(daysNmb, d)
		if iDay != -1 {
			if d == "-2" {
				daysNmb = slices.Replace(daysNmb, iDay, iDay+1, fmt.Sprint(lastMonthday(nxt)-1))
			} else {
				daysNmb = slices.Replace(daysNmb, iDay, iDay+1, fmt.Sprint(lastMonthday(nxt)))
			}
		}
	}

	return daysNmb
}

// вычисляет последнюю дату месяца по текущему времени
func lastMonthday(now time.Time) int {
	max := now
	for max.Day() > 1 {
		max = max.AddDate(0, 0, 1)
	}
	max = max.AddDate(0, 0, -1)

	return max.Day()
}
