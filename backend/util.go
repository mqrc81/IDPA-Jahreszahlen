// General helper functions to be used throughout the project

package backend

import (
	"fmt"
	"log"
	"regexp"
	"time"
)

// Abs gets absolute value of an int number (-10 => 10)
func Abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}

// Regex checks if a certain regular expression matches a certain string.
func Regex(str string, regex string) bool {
	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Fatal(fmt.Errorf("error comparing regular-expression to string: %w", err))
	}

	return match
}

// Date creates a new date from year, month and day
func Date(year int, month int, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}
