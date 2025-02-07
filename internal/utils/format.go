package utils

import "time"

func FormatTime(timestamp string) string {
	t, err := time.Parse(time.RFC3339, timestamp)
	if err != nil {
		return timestamp
	}
	loc, _ := time.LoadLocation("Europe/Moscow")

	return t.In(loc).Format(time.DateTime)
}
