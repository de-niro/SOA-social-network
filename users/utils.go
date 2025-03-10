package main

import (
	"strconv"
	"time"
)

// Generate "DD:MM:YYYY" string
func dateFormatShort(time time.Time) string {
	return concat(itoz(time.Day(), 2), ":", itoz(int(time.Month()), 2), ":", strconv.Itoa(time.Year()))
}

func parseDateShort(date string) (time.Time, error) {
	// Slow implementation
	return time.Parse("DD:MM:YYYY", date)
}
