package model

import (
	"errors"
	"slices"
	"strconv"
	"strings"
	"time"
)

const (
	dayInWeek   = 7
	monthInYear = 12
)

// NextDate - calculates the next date according to the specified rule
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// start - the date of Task start
	var (
		start time.Time // start task time
		next  time.Time // next task time (after curret)
		err   error
	)

	if err := checkRepeatFormat(repeat); err != nil {
		return "", err
	}

	start, err = time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch strings.Split(repeat, " ")[0] {
	case "y":
		// сразу прибавляем год
		next = start.AddDate(1, 0, 0)

		// если прибавленного года не хватает
		for now.Compare(next) >= 0 {
			next = next.AddDate(1, 0, 0)
		}

		return next.Format("20060102"), nil

	case "d":
		interval, err := strconv.Atoi(strings.Split(repeat, " ")[1])
		if err != nil {
			return "", err
		}

		next = start.AddDate(0, 0, interval)

		for now.Compare(next) >= 0 {
			next = start.AddDate(0, 0, interval)
		}

		return next.Format("20060102"), nil

	case "w":
		// получить дни недели скорректированные на дни от 1 до 0
		daysWeek := correctWeek(strings.Split(strings.Split(repeat, " ")[1], ","))

		// создать календарь текущей недели и отметить дни недели
		ruleClndr, err := ruleCalendar(createClndr(dayInWeek), daysWeek)
		if err != nil {
			return "", err
		}

		// найти следующий день задачи после стартового
		next = start.AddDate(0, 0, 1)

		// день недели старта задачи
		for !ruleClndr[int(next.Weekday())] || (now.Compare(next) >= 0) {
			next = next.AddDate(0, 0, 1)
		}

		return next.Format("20060102"), nil

	case "m":
		rule := strings.Split(repeat, " ")

		// create a monthday rule slice
		month := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
		if len(rule) == 3 {
			month = strings.Split(rule[2], ",")
		}

		// create calendar map and apply a repetition rule
		monthsClndr := createClndr(monthInYear)
		ruleMonthsClndr, err := ruleCalendar(monthsClndr, month)
		if err != nil {
			return "", err
		}

		// next должен быть больше start хотябы на один день
		next := start.AddDate(0, 0, 1)

		for {
			// next принадлежит текущему месяцу ?
			if ruleMonthsClndr[int(next.Month())] {
				// создаем кадендарь дней и отмечаем эти дни в текущем месяце
				ruleDaysClndr, err := calendarMonthDay(next, rule)
				if err != nil {
					return "", err
				}

				// итерации по календарю месяца
				for !ruleDaysClndr[next.Day()] { // если нет очередного дня или месяца
					next = next.AddDate(0, 0, 1)
				}

				// найденная дата до текущей даты?
				if now.Before(next) {
					return next.Format("20060102"), nil
				}

			} else {
				next = next.AddDate(0, 0, 1)
			}
		}

	}
	return "", errors.New("invalid character of repeat mod")

}

func calendarMonthDay(nextDate time.Time, rule []string) (map[int]bool, error) {
	maxDays := daysInMonth(nextDate)
	daysClndr := createClndr(maxDays)
	ruleDays := strings.Split(rule[1], ",")
	// replace "-1", "-2"
	ruleDays = replaceRelativeDates(ruleDays, maxDays)

	return ruleCalendar(daysClndr, ruleDays)
}

// вычисляет последнюю дату месяца по текущему времени
func daysInMonth(now time.Time) int {
	max := now
	for max.Day() > 1 {
		max = max.AddDate(0, 0, 1)
	}
	max = max.AddDate(0, 0, -1)

	return max.Day()
}

// заменяет в слайсе относительные значения дней месяца "-1", "-2"
func replaceRelativeDates(days []string, maxDay int) []string {
	last := strconv.Itoa(maxDay)
	penultimate := strconv.Itoa(maxDay - 1)

	if slices.Contains(days, "-1") {
		days[slices.Index(days, "-1")] = last
	}

	if slices.Contains(days, "-2") {
		days[slices.Index(days, "-2")] = penultimate
	}
	return days
}

// создает map с днями недели или месяца
func createClndr(maxItems int) map[int]bool {
	dates := map[int]bool{}
	for i := 1; i <= maxItems; i++ {
		dates[i] = false
	}
	return dates
}

func ruleCalendar(calendar map[int]bool, dates []string) (map[int]bool, error) {
	numDates, err := charsToInts(dates)
	if err != nil {
		return map[int]bool{}, err
	}

	for _, numDate := range numDates {
		calendar[numDate] = true
	}

	return calendar, nil
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

func correctWeek(days []string) []string {
	var day []string

	for _, d := range days {
		if d == "7" {
			d = "0"
		}
		day = append(day, d)
	}
	slices.Sort(day)
	return day
}
