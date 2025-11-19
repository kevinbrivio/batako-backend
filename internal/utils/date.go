package utils

import (
	"time"
)

func GetDayRange(now time.Time) (time.Time, time.Time) {
	startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// End of day: 23:59:59.999999 on the same day
	endOfDay := startOfDay.Add(23*time.Hour + 59*time.Minute + 59*time.Second + 999*time.Millisecond + 999*time.Microsecond)

	return startOfDay, endOfDay
}

func GetWeekRange(now time.Time, weekOffset int) (time.Time, time.Time) {
	// find this week's monday
	weekday := int(now.Weekday())
	if weekday == 0 { // Sunday is classified as 0 in Go
		weekday = 7
	}

	startOfWeek := now.AddDate(0, 0, -(weekday-1)-(weekOffset*7))
	startOfWeek = time.Date(startOfWeek.Year(), startOfWeek.Month(), startOfWeek.Day(), 0, 0, 0, 0, startOfWeek.Location())

	endOfWeek := startOfWeek.AddDate(0, 0, 6).Add(time.Hour*23 + time.Minute*59 + time.Second*59 + time.Millisecond*999 + time.Microsecond*999)

	return startOfWeek, endOfWeek
}

func GetMonthRange(now time.Time, monthOffset int) (time.Time, time.Time) {
	year, month, _ := now.Date()
	targetMonth := time.Month(month) + time.Month(monthOffset)
	targetYear := year

	if targetMonth < time.January {
		targetMonth += 12
		targetYear--
	} else if targetMonth > time.December {
		targetMonth -= 12
		targetYear++
	}

	startOfMonth := time.Date(targetYear, targetMonth, 1, 0, 0, 0, 0, now.UTC().Location())

	// End of target month (first of next month, then day=0)
	firstOfNextMonth := startOfMonth.AddDate(0, 1, 0)
	endOfMonth := time.Date(firstOfNextMonth.Year(), firstOfNextMonth.Month(), 0, 23, 59, 59, 0, now.UTC().Location())

	return startOfMonth, endOfMonth
}
