package helper

import "time"

func GetStartOfDay(t time.Time) time.Time {
	tUTC := t.UTC()
	return time.Date(tUTC.Year(), tUTC.Month(), tUTC.Day(), 0, 0, 0, 0, time.UTC)
}

func GetEndOfDay(t time.Time) time.Time {
	tUTC := t.UTC()
	return time.Date(tUTC.Year(), tUTC.Month(), tUTC.Day(), 23, 59, 59, 999999999, time.UTC)
}

func CalculateWorkingHours(checkIn, checkOut time.Time, lunchDuration int) float64 {

	duration := checkOut.Sub(checkIn)

	hours := duration.Hours() - float64(lunchDuration)/60.0
	if hours < 0 {
		return 0
	}	

	return hours

}
