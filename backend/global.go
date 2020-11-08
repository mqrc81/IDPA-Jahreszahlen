package backend

// "Thema"
type Unit struct {
	ID          int    `db:"id"`
	StartYear   int    `db:"start_year"`
	EndYear     int    `db:"end_year"`
	Title       string `db:"title"`
	Description string `db:"description"`
	PlayCount   int    `db:"playcount"`
}

// "Ereignis"
type Event struct {
	ID     int    `db:"id"`
	UnitID int    `db:"unit_id"`
	Title  string `db:"title"`
	Year   int    `db:"year"`
}

type UnitStore interface {
	Unit(id int) (Unit, error)
	Units() ([]Unit, error)
	CreateUnit(u *Unit) error
	UpdateUnit(u *Unit) error
	DeleteUnit(id int) error
}

type EventStore interface {
	Event(id int) (Event, error)
	EventsByUnit(unitID int) ([]Event, error)
	CreateEvent(e *Event) error
	UpdateEvent(e *Event) error
	DeleteEvent(id int) error
}

type Store interface {
	UnitStore
	EventStore
}
