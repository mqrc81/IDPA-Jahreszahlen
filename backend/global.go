package backend

/*
 * TODO Header
 */

/*
 * Topic (ger. "Thema") marks a historical segment consisting of multiple events (e.g. World War 1)
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
 * Event (ger. "Ereignis") marks a historical event associated with a specific year (e.g. Battle of Britain)
 */
type Event struct {
	EventID int    `db:"event_id"`
	TopicID int    `db:"topic_id"`
	Title   string `db:"title"`
	Year    int    `db:"year"`
}

/*
 * TopicStore stores functions for Topic to inherit
 */
type TopicStore interface {
	Topic(topicID int) (Topic, error)
	Topics() ([]Topic, error)
	CreateTopic(u *Topic) error
	UpdateTopic(u *Topic) error
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
 * Store inherits TopicStore and EventStore
 */
type Store interface {
	TopicStore
	EventStore
}
