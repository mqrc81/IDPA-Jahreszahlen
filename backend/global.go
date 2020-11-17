package backend

/*
 * TODO Header
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
 * User represents an account created
 */
type User struct {
	Username string `db:"username"`
	Password string `db:"password"`
	Admin    bool   `db:"admin"`
}

/*
 * Score represents points scored by a user upon finishing a topic
 */
type Score struct {
	ScoreID  int    `db:"score_id"`
	TopicID  int    `db:"topic_id"`
	Username string `db:"username"`
	Points   int    `db:"points"`
	Date     string `db:"date"`
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
	User(username string) (User, error)
	CreateUser(u *User) error
	UpdateUser(u *User) error
	DeleteUser(username string) error
}

/*
 * ScoreStore stores functions for Score to inherit
 */
type ScoreStore interface {
	ScoresByTopic(topicID int, limit int) ([]Score, error)
	ScoresByUser(username string, limit int) ([]Score, error)
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
