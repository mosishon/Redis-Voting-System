package helpers

import "time"

func InTimeSpan(startDate, endDate time.Time, now time.Time) bool {
	return now.After(startDate) && now.Before(endDate)
}

func RFC3339CurrentTime() time.Time {
	time, _ := time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
	return time
}
