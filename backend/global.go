package backend

/*
 * global.go contains declaration of all global structs and interfaces to be used throughout the project
 */

/*
 * Topic (ger. "Thema") represents a historical segment consisting of multiple events (e.g. World War 1)
 */
type Topic struct {
	TopicID     int    `db:"topic_id"`
	Title       string `db:"title"`
	StartYear   int    `db:"start_year"`
	EndYear     int    `db:"end_year"`
	Description string `db:"description"`
	PlayCount   int    `db:"playcount"`
}

/*
 * Event (ger. "Ereignis") represents a historical event associated with a specific year (e.g. Battle of Britain)
 */
type Event struct {
	EventID int    `db:"event_id"`
	TopicID int    `db:"topic_id"`
	Title   string `db:"title"`
	Year    int    `db:"year"`
}

/*
 * User (ger. "Benutzer") represents an account created
 */
type User struct {
	UserID   int    `db:"user_id"`
	Username string `db:"username"`
	Password string `db:"password"`
	Admin    bool   `db:"admin"`
}

/*
 * Score (ger. "Resultat") represents points scored by a user upon finishing playing a topic
 */
type Score struct {
	ScoreID int    `db:"score_id"`
	TopicID int    `db:"topic_id"`
	UserID  int    `db:"user_id"`
	Points  int    `db:"points"`
	Date    string `db:"date"`
}

/*
 * TopicStore stores functions for Topic to inherit
 */
type TopicStore interface {
	Topic(topicID int) (Topic, error)
	Topics() ([]Topic, error)
	CreateTopic(t *Topic) error
	UpdateTopic(t *Topic) error
	DeleteTopic(topicID int) error
}

/*
 * EventStore stores functions for Event to inherit
 */
type EventStore interface {
	Event(eventID int) (Event, error)
	EventsByTopic(topicID int, orderByRand bool) ([]Event, error)
	CreateEvent(e *Event) error
	UpdateEvent(e *Event) error
	DeleteEvent(eventID int) error
}

/*
 * UserStore stores functions for User to inherit
 */
type UserStore interface {
	User(userID int) (User, error)
	UserByUsername(username string) (User, error)
	Users() ([]User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(userID int) error
}

/*
 * ScoreStore stores functions for Score to inherit
 */
type ScoreStore interface {
	Scores(limit int, offset int) ([]Score, error)
	ScoresByTopic(topicID int, limit int, offset int) ([]Score, error)
	ScoresByUser(userID int, limit int, offset int) ([]Score, error)
	ScoresByTopicAndUser(topicID int, userID int, limit int, offset int) ([]Score, error)
	CreateScore(s *Score) error
}

/*
 * Store inherits TopicStore, EventStore, UserStore and ScoreStore
 */
type Store interface {
	TopicStore
	EventStore
	UserStore
	ScoreStore
}
