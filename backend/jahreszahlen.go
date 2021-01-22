// The collection of all global structures and interfaces to be used throughout
// the application. The structures are not equivalent to the database tables,
// as 1:n relationships may be stored directly in the structure of the primary
// object for the sake of practicality (e.g topic containing array of events).

package backend

import (
	"time"
)

// Topic represents a historical segment consisting of multiple events.
type Topic struct {
	TopicID     int     `db:"topic_id"`
	Name        string  `db:"name"`
	StartYear   int     `db:"start_year"`
	EndYear     int     `db:"end_year"`
	Description string  `db:"description"`
	Events      []Event `db:"events"`
	ScoresCount int     `db:"scores_count"`
	EventsCount int     `db:"events_count"`
}

// Event represents a historical event associated with a specific year.
type Event struct {
	EventID int       `db:"event_id"`
	TopicID int       `db:"topic_id"`
	Name    string    `db:"name"`
	Year    int       `db:"year"`
	Date    time.Time `db:"date"` // only relevant for sorting the events by date, if 2 events are in the same year
}

// User represents a person's account.
type User struct {
	UserID      int    `db:"user_id"`
	Username    string `db:"username"`
	Email       string `db:"email"`
	Password    string `db:"password"`
	Admin       bool   `db:"admin"`
	Verified    bool   `db:"verified"`
	ScoresCount int    `db:"scores_count"`
}

// Score represents points scored by a user upon having successfully finished
// playing a quiz.
type Score struct {
	ScoreID   int       `db:"score_id"`
	TopicID   int       `db:"topic_id"`
	UserID    int       `db:"user_id"`
	Points    int       `db:"points"`
	Date      time.Time `db:"date"`
	TopicName string    `db:"topic_name"`
	UserName  string    `db:"user_name"`
}

// Token represents a token to be sent to the user by email in case of a
// forgotten password.
type Token struct {
	TokenID string    `db:"token_id"`
	UserID  int       `db:"user_id"`
	Expiry  time.Time `db:"expiry"`
}

// TopicStore stores functions using topics for the database-layer.
type TopicStore interface {
	GetTopic(topicID int) (Topic, error)
	GetTopics() ([]Topic, error)
	CountTopics() (int, error) // unused
	CreateTopic(topic *Topic) error
	UpdateTopic(topic *Topic) error
	DeleteTopic(topicID int) error
}

// EventStore stores functions using events for the database-layer.
type EventStore interface {
	GetEvent(eventID int) (Event, error)
	CountEvents() (int, error) // unused
	CreateEvent(event *Event) error
	UpdateEvent(event *Event) error
	DeleteEvent(eventID int) error
}

// UserStore stores functions using users for the database-layer.
type UserStore interface {
	GetUser(userID int) (User, error)
	GetUserByUsername(username string) (User, error)
	GetUserByEmail(email string) (User, error)
	GetUsers() ([]User, error)
	GetAdmins() ([]User, error)
	CountUsers() (int, error) // unused
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(userID int) error
}

// ScoreStore stores functions using scores for the database-layer.
type ScoreStore interface {
	GetScores() ([]Score, error)
	GetScoresByTopic(topicID int) ([]Score, error)
	GetScoresByUser(userID int) ([]Score, error)
	GetScoresByTopicAndUser(topicID int, userID int) ([]Score, error)
	CreateScore(score *Score) error
}

// TokenStore stores functions using tokens for the database-layer.
type TokenStore interface {
	GetToken(tokenID string) (Token, error)
	CreateToken(token *Token) error
	DeleteTokensByUser(userID int) error
}

// Store combines TopicStore, EventStore, UserStore and ScoreStore.
type Store interface {
	TopicStore
	EventStore
	UserStore
	ScoreStore
	TokenStore
}
