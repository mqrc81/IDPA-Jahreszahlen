// A collection of general helper functions.

package web

import (
	"log"
	"regexp"
)

// min returns the smallest out of all the numbers.
func min(nums ...int) int {

	if len(nums) == 0 {
		return 0
	}

	minNumber := nums[0]
	for _, num := range nums {
		if num < minNumber {
			minNumber = num
		}
	}

	return minNumber
}

// abs returns the absolute value of a number.
func abs(num int) int {

	if num < 0 {
		return -num
	}

	return num
}

// url catches empty URLs, which my occur when manually typing in a URL to
// which a user doesn't have access to, in which case the 'req.Referer()' would
// be empty. In case of an empty URL it redirects to the home-page.
func url(url string) string {

	if url == "" {
		return "/"
	}

	return url
}

// regex checks if a certain regular expression matches a certain string.
func regex(str string, regex string) bool {

	match, err := regexp.MatchString(regex, str)
	if err != nil {
		log.Printf("error comparing regular-expression '%v' to string '%v': %v", regex, str, err)
	}

	return match
}
