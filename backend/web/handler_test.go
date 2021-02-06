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

// TestScoreCreatePages (from score_handler) tests creating the pages a user can
// navigate to from the leaderboard ('[< 1 '2' 3 >]').
func TestScoreCreatePages(t *testing.T) {

	// Declare test cases
	tests := []struct {
		name        string
		page        int
		show        int
		scoresCount int
		want        []int
	}{
		{
			name:        "#1 NORMAL",
			page:        3,
			show:        10,
			scoresCount: 45,
			want:        []int{2, 3, 4},
		},
		{
			name:        "#2 FIRST PAGE",
			page:        1,
			show:        10,
			scoresCount: 45,
			want:        []int{1, 2, 3},
		},
		{
			name:        "#3 LAST PAGE",
			page:        5,
			show:        10,
			scoresCount: 45,
			want:        []int{3, 4, 5},
		},
		{
			name:        "#4 FIRST OF 2 PAGES",
			page:        1,
			show:        10,
			scoresCount: 15,
			want:        []int{1, 2},
		},
		{
			name:        "#5 LAST OF 2 PAGES",
			page:        2,
			show:        10,
			scoresCount: 15,
			want:        []int{1, 2},
		},
		{
			name:        "#6 1 PAGE",
			page:        1,
			show:        10,
			scoresCount: 5,
			want:        []int{1},
		},
		{
			name:        "#7 LAST ROW IS LAST SCORE",
			page:        3,
			show:        10,
			scoresCount: 30,
			want:        []int{1, 2, 3},
		},
		{
			name:        "#8 LAST ROW IS LAST SCORE MINUS 1",
			page:        3,
			show:        10,
			scoresCount: 29,
			want:        []int{1, 2, 3},
		},
		{
			name:        "#9 LAST ROW IS LAST SCORE PLUS 1",
			page:        3,
			show:        10,
			scoresCount: 31,
			want:        []int{2, 3, 4},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := createPages(test.page, test.show, test.scoresCount); !reflect.DeepEqual(got, test.want) {
				t.Errorf("createPages() = %v, want %v", got, test.want)
			}
		})
	}
}
