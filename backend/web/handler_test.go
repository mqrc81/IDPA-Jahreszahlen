// Collection of tests for local private function for very specific tasks, that
// are tricky to walk through for each edge-case. HTTP-handlers or similar are
// not tested here.

package web

import (
	"reflect"
	"testing"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// tScores is a mock array of scores for testing purposes
	tScores []x.Score

	// nilScores is a nil slice of scores, since "var s []Score" is a nil slice
	// but "s := []Score{}" is an empty slice (so we can't use the latter for
	// this use case)
	nilScores []x.Score
)

// Skip other init functions in the package, which includes parsing templates,
// which would resolve in an error. Package level variables get initialized
// before the init function, thus the init function gets skipped when running
// these tests.
var _ = func() interface{} {
	_testing = true
	return nil
}()

func init() {
	// Create mock scores
	for i := 100; i > 0; i -= 3 {
		score := x.Score{
			Points: i,
		}
		tScores = append(tScores, score)
	}
}

// TestQuizBinarySearchForPoints (from quiz_handler) tests searching for the index at which the user's
// points would rank in, if all current scores from the same given topic were
// sorted by points in descending order.
func TestQuizBinarySearchForPoints(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name   string
		points int
		scores []x.Score
		want   int
	}{
		{
			name:   "#1 30 POINTS",
			points: 30,
			scores: tScores,
			want:   24,
		},
		{
			name:   "#2 HIGHEST POINTS",
			points: 1000,
			scores: tScores,
			want:   0,
		},
		{
			name:   "#3 LOWEST POINTS",
			points: 0,
			scores: tScores,
			want:   len(tScores),
		},
		{
			name:   "#4 MIDDLE POINTS",
			points: 50,
			scores: tScores,
			want:   len(tScores)/2 - 1,
		},
		{
			name:   "#5 NO SCORES",
			points: 50,
			scores: nilScores,
			want:   0,
		},
		{
			name:   "#6 1 HIGHER SCORE",
			points: 50,
			scores: tScores[:1],
			want:   1,
		},
		{
			name:   "#7 1 LOWER SCORE",
			points: 50,
			scores: tScores[len(tScores)-1:],
			want:   0,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := binarySearchForPoints(test.points, test.scores, 0, len(test.scores)); got != test.want {
				t.Errorf("binarySearchForPoints() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestScoreInspectFilters (from score_handler) tests changing the filters to
// values that are possible for the leaderboard to display.
func TestScoreInspectFilters(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name        string
		show        string
		page        string
		scoresCount int
		wantShow    int
		wantPage    int
	}{
		{
			name:        "#1 NORMAL",
			show:        "10",
			page:        "2",
			scoresCount: 25,
			wantShow:    10,
			wantPage:    2,
		},
		{
			name:        "#2 0 SCORES",
			show:        "10",
			page:        "2",
			scoresCount: 0,
			wantShow:    10,
			wantPage:    1,
		},
		{
			name:        "#3 1 SCORE",
			show:        "10",
			page:        "2",
			scoresCount: 1,
			wantShow:    10,
			wantPage:    1,
		},
		{
			name:        "#4 SAME SCORES AS SHOW",
			show:        "25",
			page:        "22",
			scoresCount: 25,
			wantShow:    25,
			wantPage:    1,
		},
		{
			name:        "#5 1 MORE SCORE THAN SHOW",
			show:        "25",
			page:        "2",
			scoresCount: 26,
			wantShow:    25,
			wantPage:    2,
		},
		{
			name:        "#6 1 FEWER SCORE THAN SHOW",
			show:        "25",
			page:        "2",
			scoresCount: 24,
			wantShow:    25,
			wantPage:    1,
		},
		{
			name:        "#7 SHOW ALL SCORES",
			show:        "-1",
			page:        "2",
			scoresCount: 420,
			wantShow:    420,
			wantPage:    1,
		},
		{
			name:        "#8 SHOW MISSING",
			page:        "3",
			scoresCount: 5,
			wantShow:    showDefault,
			wantPage:    1,
		},
		{
			name:        "#9 PAGE MISSING",
			show:        "25",
			scoresCount: 100,
			wantShow:    25,
			wantPage:    1,
		},
		{
			name:        "#10 INVALID SHOW BELOW 20",
			show:        "17",
			page:        "2",
			scoresCount: 35,
			wantShow:    10,
			wantPage:    2,
		},
		{
			name:        "#11 INVALID SHOW BELOW 40",
			show:        "33",
			page:        "2",
			scoresCount: 15,
			wantShow:    25,
			wantPage:    1,
		},
		{
			name:        "#12 INVALID SHOW ABOVE 40",
			show:        "69",
			page:        "3",
			scoresCount: 69,
			wantShow:    50,
			wantPage:    2,
		},
	}

	// Run tests
	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			show, page := inspectFilters(test.show, test.page, test.scoresCount)

			if show != test.wantShow || page != test.wantPage {
				t.Errorf("createPages() show, page = %v, %v, want %v, %v", show, page, test.wantShow, test.wantPage)
			}
		})
	}
}

// TestScoreCreatePages (from score_handler) tests creating the pages a user can
// navigate to from the leaderboard ('[< 1 '2' 3 >]').
func TestScoreCreatePages(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name        string
		show        int
		page        int
		scoresCount int
		want        []int
	}{
		{
			name:        "#1 NORMAL",
			show:        10,
			page:        3,
			scoresCount: 45,
			want:        []int{2, 3, 4},
		},
		{
			name:        "#2 FIRST PAGE",
			show:        10,
			page:        1,
			scoresCount: 45,
			want:        []int{1, 2, 3},
		},
		{
			name:        "#3 LAST PAGE",
			show:        10,
			page:        5,
			scoresCount: 45,
			want:        []int{3, 4, 5},
		},
		{
			name:        "#4 FIRST OF 2 PAGES",
			show:        10,
			page:        1,
			scoresCount: 15,
			want:        []int{1, 2},
		},
		{
			name:        "#5 LAST OF 2 PAGES",
			show:        10,
			page:        2,
			scoresCount: 15,
			want:        []int{1, 2},
		},
		{
			name:        "#6 1 PAGE",
			show:        10,
			page:        1,
			scoresCount: 5,
			want:        []int{1},
		},
		{
			name:        "#7 LAST ROW IS LAST SCORE",
			show:        10,
			page:        3,
			scoresCount: 30,
			want:        []int{1, 2, 3},
		},
		{
			name:        "#8 LAST ROW IS LAST SCORE MINUS 1",
			show:        10,
			page:        3,
			scoresCount: 29,
			want:        []int{1, 2, 3},
		},
		{
			name:        "#9 LAST ROW IS LAST SCORE PLUS 1",
			show:        10,
			page:        3,
			scoresCount: 31,
			want:        []int{2, 3, 4},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := createPages(test.show, test.page, test.scoresCount); !reflect.DeepEqual(got, test.want) {
				t.Errorf("createPages() = %v, want %v", got, test.want)
			}
		})
	}
}

// TestGenerateRandomString tests generating a random string of a certain length.
func TestGenerateRandomString(t *testing.T) {

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
			got := generateRandomString(test.len)

			if len(got) != test.wantLen || reflect.TypeOf(got) != reflect.TypeOf("") {
				t.Errorf("GenerateString() = %v, want string of length %v", got, test.wantLen)
				return
			}
		})
	}
}
