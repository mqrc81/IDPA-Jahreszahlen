// General helper functions to be used throughout the project.

package util

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"regexp"
	"time"
)

// Abs gets absolute value of an int number.
// (-10 => 10)
func Abs(num int) int {

	if num < 0 {
		return -num
	}

	return num
}

// Max returns the biggest value.
func Max(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	max := nums[0]
	for _, num := range nums {
		if num > max {
			max = num
		}
	}

	return max
}

// Min returns the smallest value.
func Min(nums ...int) int {
	if len(nums) == 0 {
		return 0
	}

	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}

	return min
}

// Regex checks if a certain regular expression matches a certain string.
func Regex(str string, regex string) bool {

	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Printf("error comparing regular-expression to string: %v", err)
	}

	return match
}

// Date creates a new date from year, month and day.
func Date(year int, month int, day int) time.Time {

	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

// GenerateBytes generates a random byte key of a certain length.
func GenerateBytes(len int) []byte {

	key := make([]byte, len)

	if _, err := rand.Read(key); err != nil {
		log.Fatalf("error generating random key: %v", err)
	}

	return key
}

// GenerateString generates a random string key of string of a certain length.
func GenerateString(len int) string {

	key := GenerateBytes(len)

	return base64.URLEncoding.EncodeToString(key)[:len]
}
