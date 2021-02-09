// Collection of tests for the database access layer of functions evolving
// around scores.

package database

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// tScore is a mock score for testing purposes
	tScore = x.Score{
		ScoreID:   1,
		TopicID:   1,
		UserID:    1,
		Points:    50,
		Date:      time.Now(),
		TopicName: "Topic 1",
		UserName:  "user_1",
	}

	// tScore2 is a mock score for testing purposes
	tScore2 = x.Score{
		ScoreID:   2,
		TopicID:   1,
		UserID:    2,
		Points:    60,
		Date:      time.Now().Add(time.Hour * 1),
		TopicName: "Topic 1",
		UserName:  "user_2",
	}

	// tScore3 is a mock score for testing purposes
	tScore3 = x.Score{
		ScoreID:   3,
		TopicID:   2,
		UserID:    1,
		Points:    70,
		Date:      time.Now().Add(time.Hour * 2),
		TopicName: "Topic 2",
		UserName:  "user_1",
	}

	// nilScores is a nil slice of scores, since "var s []Score" is a nil slice
	// but "s := []Score{}" is an empty slice (so we can't use the latter for
	// this use case)
	nilScores []x.Score
)

// TestGetScores tests getting all scores.
func TestGetScores(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM scores"

	tScores := []x.Score{tScore, tScore2, tScore3}
	table := []string{"score_id", "topic_id", "user_id", "points", "date", "topic_name", "user_name"}

	// Declare test cases
	tests := []struct {
		name       string
		mock       func()
		wantScores []x.Score
		wantError  bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows(table)
				for _, score := range tScores {
					rows = rows.AddRow(score.ScoreID, score.TopicID, score.UserID, score.Points, score.Date,
						score.TopicName, score.UserName)
				}

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantScores: tScores,
			wantError:  false,
		},
		{
			// When the scores table is empty
			name: "#2 OK (NO ROWS)",
			mock: func() {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantScores: nilScores,
			wantError:  false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock()

			scores, err := store.GetScores()

			if (err != nil) != test.wantError {
				t.Errorf("GetScores() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scores, test.wantScores) {
				t.Errorf("GetScores() = %v, want %v", scores, test.wantScores)
			}
		})
	}
}

// TestGetScoresByTopic tests getting all scores of a certain topic.
func TestGetScoresByTopic(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM scores"

	tScores := []x.Score{tScore, tScore2}
	table := []string{"score_id", "topic_id", "user_id", "points", "date", "topic_name", "user_name"}

	// Declare test cases
	tests := []struct {
		name       string
		topicID    int
		mock       func(topicID int)
		wantScores []x.Score
		wantError  bool
	}{
		{
			// When everything works as intended
			name:    "#1 OK",
			topicID: tScores[0].TopicID,
			mock: func(topicID int) {
				rows := sqlmock.NewRows(table)
				for _, score := range tScores {
					rows = rows.AddRow(score.ScoreID, score.TopicID, score.UserID, score.Points, score.Date,
						score.TopicName, score.UserName)
				}

				mock.ExpectQuery(queryMatch).WithArgs(topicID).WillReturnRows(rows)
			},
			wantScores: tScores,
			wantError:  false,
		},
		{
			// When the scores table is empty
			name:    "#2 OK (NO ROWS)",
			topicID: tScores[0].TopicID,
			mock: func(topicID int) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(topicID).WillReturnRows(rows)
			},
			wantScores: nilScores,
			wantError:  false,
		},
		{
			// When topic with given topic ID doesn't exist
			name:    "#3 TOPIC NOT FOUND",
			topicID: 0,
			mock: func(topicID int) {
				mock.ExpectQuery(queryMatch).WithArgs(topicID).
					WillReturnError(errors.New("topic with given id does not exist"))
			},
			wantScores: nilScores,
			wantError:  true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.topicID)

			scores, err := store.GetScoresByTopic(test.topicID)

			if (err != nil) != test.wantError {
				t.Errorf("GetScoresByTopic() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scores, test.wantScores) {
				t.Errorf("GetScoresByTopic() = %v, want %v", scores, test.wantScores)
			}
		})
	}
}

// TestGetScoresByTopicAndUser tests getting all scores of a certain topic and
// a certain user.
func TestGetScoresByTopicAndUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM scores"

	tScores := []x.Score{tScore}
	table := []string{"score_id", "topic_id", "user_id", "points", "date", "topic_name", "user_name"}

	// Declare test cases
	tests := []struct {
		name       string
		topicID    int
		userID     int
		mock       func(topicID int, userID int)
		wantScores []x.Score
		wantError  bool
	}{
		{
			// When everything works as intended
			name:    "#1 OK",
			topicID: tScores[0].TopicID,
			userID:  tScores[0].UserID,
			mock: func(topicID int, userID int) {
				rows := sqlmock.NewRows(table)
				for _, score := range tScores {
					rows = rows.AddRow(score.ScoreID, score.TopicID, score.UserID, score.Points, score.Date,
						score.TopicName, score.UserName)
				}

				mock.ExpectQuery(queryMatch).WithArgs(topicID, userID).WillReturnRows(rows)
			},
			wantScores: tScores,
			wantError:  false,
		},
		{
			// When the scores table is empty
			name:    "#2 OK (NO ROWS)",
			topicID: tScores[0].TopicID,
			userID:  tScores[0].UserID,
			mock: func(topicID int, userID int) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(topicID, userID).WillReturnRows(rows)
			},
			wantScores: nilScores,
			wantError:  false,
		},
		{
			// When topic with given topic ID doesn't exist
			name:    "#3 TOPIC NOT FOUND",
			topicID: 0,
			userID:  tScores[0].UserID,
			mock: func(topicID int, userID int) {
				mock.ExpectQuery(queryMatch).WithArgs(topicID, userID).
					WillReturnError(errors.New("topic with given id does not exist"))
			},
			wantScores: nilScores,
			wantError:  true,
		},
		{
			// When user with given user ID doesn't exist
			name:    "#4 USER NOT FOUND",
			topicID: tScores[0].TopicID,
			userID:  0,
			mock: func(topicID int, userID int) {
				mock.ExpectQuery(queryMatch).WithArgs(topicID, userID).
					WillReturnError(errors.New("user with given id does not exist"))
			},
			wantScores: nilScores,
			wantError:  true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.topicID, test.userID)

			scores, err := store.GetScoresByTopicAndUser(test.topicID, test.userID)

			if (err != nil) != test.wantError {
				t.Errorf("GetScoresByTopicAndUser() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scores, test.wantScores) {
				t.Errorf("GetScoresByTopicAndUser() = %v, want %v", scores, test.wantScores)
			}
		})
	}
}

// TestCountScores tests getting amount of scores.
func TestCountScores(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT COUNT((.+)) FROM scores"

	table := []string{"COUNT(*)"}

	// Declare test cases
	tests := []struct {
		name            string
		mock            func()
		wantScoresCount int
		wantError       bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows(table).AddRow(3)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantScoresCount: 3,
			wantError:       false,
		},
		{
			// When the scores table is empty
			name: "#2 NO ROWS",
			mock: func() {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows).WillReturnError(errors.New("no scores found"))
			},
			wantScoresCount: 0,
			wantError:       true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock()

			scoresCount, err := store.CountScores()

			if (err != nil) != test.wantError {
				t.Errorf("CountScores() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scoresCount, test.wantScoresCount) {
				t.Errorf("CountScores() = %v, want %v", scoresCount, test.wantScoresCount)
			}
		})
	}
}

// TestCountScoresByDate tests getting amount of scores by date.
func TestCountScoresByDate(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT COUNT((.+)) FROM scores"

	table := []string{"COUNT(*)"}

	// Declare test cases
	tests := []struct {
		name            string
		start           time.Time
		end             time.Time
		mock            func(start time.Time, to time.Time)
		wantScoresCount int
		wantError       bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			start: time.Now().AddDate(0, -1, 0),
			end:   time.Now().AddDate(0, 0, 0),
			mock: func(start time.Time, end time.Time) {
				rows := sqlmock.NewRows(table).AddRow(3)

				mock.ExpectQuery(queryMatch).WithArgs(start, end).WillReturnRows(rows)
			},
			wantScoresCount: 3,
			wantError:       false,
		},
		{
			// When the scores table is empty
			name:  "#2 NO ROWS",
			start: time.Now().AddDate(0, -1, 0),
			end:   time.Now().AddDate(0, 0, 0),
			mock: func(start time.Time, end time.Time) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(start,
					end).WillReturnRows(rows).WillReturnError(errors.New("no scores found"))
			},
			wantScoresCount: 0,
			wantError:       true,
		},
		{
			// When the start-date is after the to-date table is empty
			name:  "#3 START AFTER END",
			start: time.Now().AddDate(0, 0, 0),
			end:   time.Now().AddDate(0, -1, 0),
			mock: func(start time.Time, end time.Time) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(start,
					end).WillReturnRows(rows).WillReturnError(errors.New("no scores found"))
			},
			wantScoresCount: 0,
			wantError:       true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.start, test.end)

			scoresCount, err := store.CountScoresByDate(test.start, test.end)

			if (err != nil) != test.wantError {
				t.Errorf("CountScoresByDate() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scoresCount, test.wantScoresCount) {
				t.Errorf("CountScoresByDate() = %v, want %v", scoresCount, test.wantScoresCount)
			}
		})
	}
}

// TestCreateScore tests creating a new score
func TestCreateScore(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "INSERT INTO scores"

	// Declare test cases
	tests := []struct {
		name      string
		score     x.Score
		mock      func(score x.Score)
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			score: tScore,
			mock: func(score x.Score) {
				mock.ExpectExec(queryMatch).WithArgs(score.TopicID, score.UserID, score.Points, score.Date).
					WillReturnResult(sqlmock.NewResult(int64(score.ScoreID), 1))
			},
			wantError: false,
		},
		{
			// When topic with given topic ID doesn't exist
			name: "#2 TOPIC NOT FOUND",
			score: x.Score{
				ScoreID: tScore.ScoreID,
				TopicID: 0,
				UserID:  tScore.UserID,
				Points:  tScore.Points,
				Date:    tScore.Date,
			},
			mock: func(score x.Score) {
				mock.ExpectExec(queryMatch).WithArgs(score.TopicID, score.UserID, score.Points, score.Date).
					WillReturnError(errors.New("topic with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When user with given user ID doesn't exist
			name: "#3 USER NOT FOUND",
			score: x.Score{
				ScoreID: tScore.ScoreID,
				TopicID: tScore.TopicID,
				UserID:  0,
				Points:  tScore.Points,
				Date:    tScore.Date,
			},
			mock: func(score x.Score) {
				mock.ExpectExec(queryMatch).WithArgs(score.TopicID, score.UserID, score.Points, score.Date).
					WillReturnError(errors.New("user with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When points are missing
			name: "#4 POINTS MISSING",
			score: x.Score{
				ScoreID: tScore.ScoreID,
				TopicID: tScore.TopicID,
				UserID:  tScore.UserID,
				Date:    tScore.Date,
			},
			mock: func(score x.Score) {
				mock.ExpectExec(queryMatch).WithArgs(score.TopicID, score.UserID, score.Points, score.Date).
					WillReturnError(errors.New("points can not be empty"))
			},
			wantError: true,
		},
		{
			// When date is missing
			name: "#5 DATE MISSING",
			score: x.Score{
				ScoreID: tScore.ScoreID,
				TopicID: tScore.TopicID,
				UserID:  tScore.UserID,
				Points:  tScore.Points,
			},
			mock: func(score x.Score) {
				mock.ExpectExec(queryMatch).WithArgs(score.TopicID, score.UserID, score.Points, score.Date).
					WillReturnError(errors.New("date can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.score)

			err := store.CreateScore(&test.score)

			if (err != nil) != test.wantError {
				t.Errorf("CreateScore() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}
