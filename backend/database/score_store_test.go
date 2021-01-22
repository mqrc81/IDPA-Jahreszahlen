package database

import (
	"testing"
	"time"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// Mock score for testing purposes
	ss = x.Score{
		ScoreID:   1,
		TopicID:   1,
		UserID:    1,
		Points:    50,
		Date:      time.Now(),
		TopicName: "Topic 1",
		UserName:  "user_1",
	}
)

// TestGetScores tests getting all scores.
func TestGetScores(t *testing.T) {

}

// TestGetScoresByTopic tests getting all scores of a certain topic.
func TestGetScoresByTopic(t *testing.T) {

}

// TestGetScoresByUser tests getting all scores of a certain user
func TestGetScoresByUser(t *testing.T) {

}

// TestGetScoresByTopicAndUser tests getting all scores of a certain topic and
// a certain user.
func TestGetScoresByTopicAndUser(t *testing.T) {

}

// TestCreateScore tests creating a new score
func TestCreateScore(t *testing.T) {

}
