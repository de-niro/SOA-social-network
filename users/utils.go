package main

import (
	"strconv"
	"time"
)

func RegistrationDateFormat(time time.Time) string {
	return concat(itoz(time.Day(), 2), ":", itoz(int(time.Month()), 2), ":", strconv.Itoa(time.Year()))
}
