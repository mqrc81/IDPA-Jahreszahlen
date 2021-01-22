package database

import (
	"errors"
	"reflect"
	"testing"
	"time"

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
	// and "s := []Score" is an empty slice
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
				rows := mock.NewRows(table)
				for _, s := range tScores {
					rows = rows.AddRow(s.ScoreID, s.TopicID, s.UserID, s.Points, s.Date, s.TopicName, s.UserName)
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
				rows := mock.NewRows(table)

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
				rows := mock.NewRows(table)
				for _, s := range tScores {
					rows = rows.AddRow(s.ScoreID, s.TopicID, s.UserID, s.Points, s.Date, s.TopicName, s.UserName)
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
				rows := mock.NewRows(table)

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

// TestGetScoresByUser tests getting all scores of a certain user
func TestGetScoresByUser(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &ScoreStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM scores"

	tScores := []x.Score{tScore, tScore3}
	table := []string{"score_id", "topic_id", "user_id", "points", "date", "topic_name", "user_name"}

	// Declare test cases
	tests := []struct {
		name       string
		userID     int
		mock       func(userID int)
		wantScores []x.Score
		wantError  bool
	}{
		{
			// When everything works as intended
			name:   "#1 OK",
			userID: tScores[0].UserID,
			mock: func(userID int) {
				rows := mock.NewRows(table)
				for _, s := range tScores {
					rows = rows.AddRow(s.ScoreID, s.TopicID, s.UserID, s.Points, s.Date, s.TopicName, s.UserName)
				}

				mock.ExpectQuery(queryMatch).WithArgs(userID).WillReturnRows(rows)
			},
			wantScores: tScores,
			wantError:  false,
		},
		{
			// When the scores table is empty
			name:   "#2 OK (NO ROWS)",
			userID: tScores[0].UserID,
			mock: func(userID int) {
				rows := mock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(userID).WillReturnRows(rows)
			},
			wantScores: nilScores,
			wantError:  false,
		},
		{
			// When user with given user ID doesn't exist
			name:   "#3 USER NOT FOUND",
			userID: 0,
			mock: func(userID int) {
				mock.ExpectQuery(queryMatch).WithArgs(userID).
					WillReturnError(errors.New("user with given id does not exist"))
			},
			wantScores: nilScores,
			wantError:  true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.userID)

			scores, err := store.GetScoresByUser(test.userID)

			if (err != nil) != test.wantError {
				t.Errorf("GetScoresByUser() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(scores, test.wantScores) {
				t.Errorf("GetScoresByUser() = %v, want %v", scores, test.wantScores)
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
				rows := mock.NewRows(table)
				for _, s := range tScores {
					rows = rows.AddRow(s.ScoreID, s.TopicID, s.UserID, s.Points, s.Date, s.TopicName, s.UserName)
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
				rows := mock.NewRows(table)

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

// TestCreateScore tests creating a new score
func TestCreateScore(t *testing.T) {

}
