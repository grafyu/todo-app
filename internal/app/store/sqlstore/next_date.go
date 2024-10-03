package sqlstore

import (
	"errors"
	"slices"
	"strings"
	"time"
)

/*
https://pkg.go.dev/slices
https://pkg.go.dev/strings
https://pkg.go.dev/time
https://goplay.space/
*/

func NextDate(now time.Time, date string, repeat string) (string, error) {
	param := strings.Split(repeat, " ")
	if !slices.Contains([]string{"y", "d", "w", "m"}, param[0]) {
		return "", errors.New("not able symbol")
	}

	if param[0] == "y" {
		start, err := time.Parse("20060102", date)
		if err != nil {
			return "", errors.New("not correct date")
		}

		next := start.AddDate(1, 0, 0).Format("20060102")

		return next, nil
	}
	return "", nil
}
