package backend

// global_vars.go
// Contains all global variables and their functions, to be accessed throughout
// the project.

// Topic
// Ger.: "Thema". Represents a historical segment consisting of multiple events
// (e.g. World War 1).
type Topic struct {
	TopicID     int    `db:"topic_id"`
	Title       string `db:"title"`
	StartYear   int    `db:"start_year"`
	EndYear     int    `db:"end_year"`
	Description string `db:"description"`
	PlayCount   int    `db:"playcount"`
}

// Event
// Ger.: "Ereignis". Represents a historical event associated with a specific
// year (e.g. Battle of Britain).
type Event struct {
	EventID int    `db:"event_id"`
	TopicID int    `db:"topic_id"`
	Title   string `db:"title"`
	Year    int    `db:"year"`
}

// User
// Ger.: "Benutzer". Represents a person's account.
type User struct {
	UserID   int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Admin    bool   `db:"admin"`
}

// Score
// Ger.: "Resultat". Represents points scored by a user upon having successfully
// finished playing a game.
type Score struct {
	ScoreID int    `db:"score_id"`
	TopicID int    `db:"topic_id"`
	UserID  int    `db:"user_id"`
	Points  int    `db:"points"`
	Date    string `db:"date"`
}

// TopicStore
// Stores functions using topics for the database-layer.
type TopicStore interface {
	Topic(topicID int) (Topic, error)
	Topics() ([]Topic, error)
	CountTopics() (int, error)
	CreateTopic(t *Topic) error
	UpdateTopic(t *Topic) error
	DeleteTopic(topicID int) error
}

// EventStore
// Stores functions using events for the database-layer.
type EventStore interface {
	Event(eventID int) (Event, error)
	EventsByTopic(topicID int, orderByRand bool) ([]Event, error)
	CountEvents() (int, error)
	CountEventsByTopic(topicID int) (int, error)
	CreateEvent(e *Event) error
	UpdateEvent(e *Event) error
	DeleteEvent(eventID int) error
}

// UserStore
// Stores functions using users for the database-layer.
type UserStore interface {
	User(userID int) (User, error)
	UserByUsername(username string) (User, error)
	Users() ([]User, error)
	CountUsers() (int, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(userID int) error
}

// ScoreStore
// Stores functions using scores for the database-layer.
type ScoreStore interface {
	Scores(limit int, offset int) ([]Score, error)
	ScoresByTopic(topicID int, limit int, offset int) ([]Score, error)
	ScoresByUser(userID int, limit int, offset int) ([]Score, error)
	ScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]Score, error)
	CountScores() (int, error)
	CreateScore(s *Score) error
}

// Store
// Holds functions of TopicStore, EventStore, UserStore and ScoreStore.
type Store interface {
	TopicStore
	EventStore
	UserStore
	ScoreStore
}
