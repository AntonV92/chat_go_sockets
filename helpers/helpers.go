package helpers

import "time"

// func to get time duration from given time until now
func GetTimeDiffNow(timeToCheck time.Time) time.Duration {
	timeToCheck = time.Date(
		timeToCheck.Year(), timeToCheck.Month(), timeToCheck.Day(),
		timeToCheck.Hour(), timeToCheck.Minute(), 0, 0, time.Local)

	return time.Now().Sub(timeToCheck)
}
