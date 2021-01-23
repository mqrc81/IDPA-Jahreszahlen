// Collection of tests for util helper functions.

package util

import (
	"reflect"
	"testing"
	"time"
)

// TestAbs tests getting the absolute value of a number.
func TestAbs(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name string
		num  int // function parameter
		want int
	}{
		{
			name: "#1 OK (POSITIVE)",
			num:  54321,
			want: 54321,
		},
		{
			name: "#2 OK (NEGATIVE)",
			num:  -123,
			want: 123,
		},
		{
			name: "#3 OK (ZERO)",
			num:  0,
			want: 0,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if got := Abs(test.num); got != test.want {
				t.Errorf("Abs() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestDate tests creating a time.Time date out of the year, month and day.
func TestDate(t *testing.T) {

	// Function parameters
	type date struct {
		year  int
		month int
		day   int
	}

	// Declare test cases
	tests := []struct {
		name string
		date date
		want time.Time
	}{
		{
			name: "#1 OK",
			date: date{
				year:  2001,
				month: 10,
				day:   12,
			},
			want: time.Date(2001, time.October, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "#2 MONTH OUT OF BOUNDS",
			date: date{
				year:  2001,
				month: 13,
				day:   12,
			},
			want: time.Date(2002, time.January, 12, 0, 0, 0, 0, time.UTC),
		},
		{
			name: "#3 DAY OUT OF BOUNDS",
			date: date{
				year:  2001,
				month: 10,
				day:   32,
			},
			want: time.Date(2001, 11, 1, 0, 0, 0, 0, time.UTC),
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if got := Date(test.date.year, test.date.month, test.date.day); !reflect.DeepEqual(got, test.want) {
				t.Errorf("Date() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestGenerateBytes tests generating a random array of bytes of a certain
// length.
func TestGenerateBytes(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name    string
		len     int // function parameter
		wantLen int
	}{
		{
			name:    "#1 OK",
			len:     32,
			wantLen: 32,
		},
		{
			name:    "#2 OK (BIG INT)",
			len:     987654321,
			wantLen: 987654321,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if got := GenerateBytes(test.len); len(got) != test.wantLen || reflect.TypeOf(got) !=
				reflect.TypeOf([]byte{}) {
				t.Errorf("GenerateBytes() = %v, want []byte of length %v", got, test.wantLen)
			}
		})
	}
}

// TestGenerateString tests generating a random string of a certain length.
func TestGenerateString(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name    string
		len     int // function parameter
		wantLen int
	}{
		{
			name:    "#1 OK",
			len:     32,
			wantLen: 32,
		},
		{
			name:    "#2 OK (BIG INT)",
			len:     987654321,
			wantLen: 987654321,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			if got := GenerateBytes(test.len); len(got) != test.wantLen || reflect.TypeOf(got) !=
				reflect.TypeOf([]byte{}) {
				t.Errorf("GenerateString() = %v, want string of length %v", got, test.wantLen)
			}
		})
	}
}

// TestRegex tests comparing a string to a regular expression.
func TestRegex(t *testing.T) {

	// Function parameters
	type compare struct {
		str   string
		regex string
	}

	// Declare test cases
	tests := []struct {
		name    string
		compare compare
		want    bool
	}{
		{
			name: "#1 OK",
			compare: compare{
				str:   "abcdefg1asd",
				regex: "\\d",
			},
			want: true,
		},
		{
			name: "#2 OK (EMAIL)",
			compare: compare{
				str:   "test@mail.com",
				regex: "^[a-z0-9._%+\\-]+@[a-z0-9.\\-]+\\.[a-z]{2,4}$",
			},
			want: true,
		},
		{
			name: "#3 INVALID ITERATION",
			compare: compare{
				str: "Passw0rd!",
				// Iterative regex checking (with ?=) is not supported in Go
				regex: "^(?=.*[a-z])(?=.*[A-Z])(?=.*\\d)(?=.*[@$!%*?&])[A-Za-z\\d@$!%*?&]{8,}$",
			},
			want: false, // error expected
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := Regex(tt.compare.str, tt.compare.regex); got != tt.want {
				t.Errorf("Regex() = %v, want %v", got, tt.want)
			}
		})
	}
}
