package backend

/*
 * Contains all global variables and their functions, to be accessed throughout
 * the project.
 */

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
	EventID int    `db:"event_id"`
	TopicID int    `db:"topic_id"`
	Name    string `db:"name"`
	Year    int    `db:"year"`
}

// User represents a person's account.
type User struct {
	UserID      int    `db:"user_id"`
	Username    string `db:"username"`
	Password    string `db:"password"`
	Admin       bool   `db:"admin"`
	ScoresCount int    `db:"scores_count"`
}

// Score represents points scored by a user upon having successfully finished
// playing a quiz.
type Score struct {
	ScoreID   int    `db:"score_id"`
	TopicID   int    `db:"topic_id"`
	UserID    int    `db:"user_id"`
	Points    int    `db:"points"`
	Date      string `db:"date"`
	TopicName string `db:"topic_name"`
	UserName  string `db:"user_name"`
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
	GetUsers() ([]User, error)
	CountUsers() (int, error) // unused
	CreateUser(user *User) error
	UpdateUser(user *User) error
	DeleteUser(userID int) error
}

// ScoreStore stores functions using scores for the database-layer.
type ScoreStore interface {
	GetScores(limit int, offset int) ([]Score, error)
	GetScoresByTopic(topicID int, limit int, offset int) ([]Score, error)
	GetScoresByUser(userID int, limit int, offset int) ([]Score, error)
	GetScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]Score, error)
	CreateScore(score *Score) error
	GetAveragePointsByTopic(topicID int) (float64, error)
}

// Store combines TopicStore, EventStore, UserStore and ScoreStore.
type Store interface {
	TopicStore
	EventStore
	UserStore
	ScoreStore
}
