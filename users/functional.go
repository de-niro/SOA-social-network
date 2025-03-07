package main

import (
	"bytes"
	"strconv"
	"strings"
)

// Helper string methods

// Efficient strings concatenation
func concat(s ...string) string {
	var sb strings.Builder
	for _, str := range s {
		sb.WriteString(str)
	}
	return sb.String()
}

// Semi-efficient zeros repeat
func repeat_zeros(times int) string {
	s := bytes.Repeat([]byte{'0'}, times)
	return string(s)
}

// Itoa function with leading zeros
func itoz(a int, padding int) string {
	// Extremely crappy implementation
	conv := strconv.Itoa(a)
	prefix := repeat_zeros(padding - len(conv))
	return concat(prefix, conv)
}
