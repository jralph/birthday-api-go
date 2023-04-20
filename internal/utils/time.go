package utils

import "time"

func IsTimeValid(t time.Time) bool {
	// Get the start of the day
	currentDay, _ := time.Parse(
		"2006-01-02",
		time.Now().Format("2006-01-02"),
	)
	future := time.Time(t).After(currentDay)
	before := time.Time(t).Before(currentDay)

	return future || (!future && !before)
}
