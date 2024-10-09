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

var (
	nextDate time.Time
	// deltaYear, deltaMonth, deltaDay int
)

// NextDate - calculates the next date according to the specified rule
func NextDate(now time.Time, date string, repeat string) (string, error) {
	// start - the date of Task start
	start, err := time.Parse("20060102", date)
	if err != nil {
		return "", err
	}

	switch strings.Split(repeat, " ")[0] {
	case "y":
		rule := strings.Split(repeat, " ")
		// check
		if len(rule) != 1 {
			return "", errors.New("not correct the repeate year format")
		}

		nextDate = start
		// calculate nextDate
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(1, 0, 0)
		}

		return nextDate.Format("20060102"), nil

	case "d":
		rule := strings.Split(repeat, " ")
		// check

		if len(rule) < 2 {
			return "", errors.New("not specified interval in days")
		} else if len(rule) > 2 {
			return "", errors.New("not correct the repeate day format")
		}

		interval, err := strconv.Atoi(rule[1])
		if err != nil {
			return "", err
		}

		if interval > 400 || interval < 1 {
			return "", errors.New("wrong number of days in an interval")
		}

		// start of next day search
		nextDate := start
		for nextDate.Before(now) {
			nextDate = nextDate.AddDate(0, 0, interval)
		}

		return nextDate.Format("20060102"), nil

	case "w":
		rule := strings.Split(repeat, " ")
		// check items
		if len(rule) < 2 {
			return "", errors.New("not specified days of the week")
		} else if len(rule) > 2 {
			return "", errors.New("not correct the repeate week format")
		}

		// check rules param
		daysWeek := strings.Split(rule[1], ",")
		numsWeek, err := charsToInts(daysWeek)
		if err != nil {
			return "", err
		}

		slices.Sort(numsWeek)
		if numsWeek[0] < 1 || numsWeek[len(numsWeek)-1] > 7 {
			return "", errors.New("wrong number of days of the week")
		}

		// create a calendar and apply a repetition rule
		clndr := createClndr(dayInWeek)
		pointWeekDays := strings.Split(rule[1], ",")
		ruleClndr, err := ruleCalendar(clndr, pointWeekDays)
		if err != nil {
			return "", err
		}

		// start of search a next_date
		nextDate := start
		for weekday := int(start.Weekday()); !ruleClndr[weekday] || nextDate.Before(now); weekday++ {
			nextDate = nextDate.AddDate(0, 0, 1)
			if weekday == 7 {
				weekday = 0
			}
		}

		return nextDate.Format("20060102"), nil

	case "m":
		rule := strings.Split(repeat, " ")

		if len(rule) < 2 {
			return "", errors.New("not specified days of month")
		}

		// check rules param
		if len(rule) == 2 {
			daysMnth := strings.Split(rule[1], ",")
			numsMnth, err := charsToInts(daysMnth)
			if err != nil {
				return "", err
			}

			slices.Sort(numsMnth)
			if numsMnth[0] < -2 || numsMnth[len(numsMnth)-1] > 31 {
				return "", errors.New("wrong number days of the month")
			}
		}

		if len(rule) == 3 {
			daysMnth := strings.Split(rule[2], ",")
			numsMnth, err := charsToInts(daysMnth)
			if err != nil {
				return "", err
			}

			slices.Sort(numsMnth)
			if numsMnth[0] < 1 || numsMnth[len(numsMnth)-1] > 12 {
				return "", errors.New("wrong number month")
			}
		}

		// create a monthday rule slice
		pointMonthYear := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12"}
		if len(rule) == 3 {
			pointMonthYear = strings.Split(rule[2], ",")
		}

		// create calendar map and apply a repetition rule
		monthsClndr := createClndr(monthInYear)
		ruleMonthsClndr, err := ruleCalendar(monthsClndr, pointMonthYear)
		if err != nil {
			return "", err
		}

		// start of search a next_date
		nextDate := start
		// внешний цикл для итерации по месяцам. в нем проверяется
		for {
			// поиск nextDay в nextMonth
			if ruleMonthsClndr[int(nextDate.Month())] && now.Before(nextDate) {
				ruleDaysClndr, err := calendarMonthDay(nextDate, rule)
				if err != nil {
					return "", err
				}

				// итерации по календарю месяца
				for !ruleDaysClndr[nextDate.Day()] || nextDate.Before(now) {
					nextDate = nextDate.AddDate(0, 0, 1)
				}
				if ruleMonthsClndr[int(nextDate.Month())] && ruleDaysClndr[nextDate.Day()] {
					return nextDate.Format("20060102"), nil
				}
			} else {
				nextDate = nextDate.AddDate(0, 0, 1)

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
