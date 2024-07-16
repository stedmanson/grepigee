package utils

import "time"

func CalculateFromTime(toTime time.Time, timeRange string) time.Time {
	switch timeRange {
	case "1h":
		return toTime.Add(-1 * time.Hour)
	case "6h":
		return toTime.Add(-6 * time.Hour)
	case "12h":
		return toTime.Add(-12 * time.Hour)
	case "1d":
		return toTime.AddDate(0, 0, -1)
	case "7d":
		return toTime.AddDate(0, 0, -7)
	default:
		return toTime.Add(-1 * time.Hour)
	}
}
