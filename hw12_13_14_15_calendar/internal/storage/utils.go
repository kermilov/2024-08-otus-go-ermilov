package storage

import "time"

const (
	Day  = time.Hour * 24
	Week = Day * 7
)

func GetStartOfDay(date time.Time) time.Time {
	return date.Truncate(Day)
}

func GetEndOfDay(date time.Time) time.Time {
	return GetStartOfDay(date).Add(Day).Add(-time.Nanosecond)
}

func GetStartOfWeek(date time.Time) time.Time {
	return date.Truncate(Week)
}

func GetEndOfWeek(date time.Time) time.Time {
	return GetStartOfWeek(date).Add(Week).Add(-time.Nanosecond)
}

func GetStartOfMonth(date time.Time) time.Time {
	return time.Date(date.Year(), date.Month(), 1, 0, 0, 0, 0, date.Location())
}

func GetEndOfMonth(date time.Time) time.Time {
	return GetStartOfMonth(date).AddDate(0, 1, 0).Add(-time.Nanosecond)
}
