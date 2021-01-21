// Collection of tests for functions accessing the database evolving around
// events.

package database

import (
	_ "database/sql"
	"fmt"
	"log"
	"reflect"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	past = time.Date(1800, time.January, 1, 0, 0, 0, 0, time.UTC)

	testEvent = x.Event{
		EventID: 1,
		TopicID: 1,
		Name:    "Test Event 1",
		Year:    1800,
		Date:    past,
	}

	testEvents = []x.Event{
		testEvent,
		{
			EventID: 2,
			TopicID: 1,
			Name:    "Test Event 2",
			Year:    1800,
			Date:    past,
		},
		{
			EventID: 3,
			TopicID: 2,
			Name:    "Test Event 3",
			Year:    1800,
			Date:    past,
		},
	}
)

// NewMock creates a new mock sqlx database.
func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	dbMock, mock, err := sqlmock.New()
	if err != nil {
		log.Fatal(fmt.Errorf("error initializing mock database: %w", err))
	}

	db := sqlx.NewDb(dbMock, "sqlmock")

	return db, mock
}

// TestGetEvent tests getting event by ID.
func TestGetEvent(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &EventStore{db}
	defer db.Close()

	// Expected result
	type want struct {
		event x.Event
		error bool
	}

	query := regexp.QuoteMeta(getEventQuery)

	// Declare test cases
	tests := []struct {
		name    string
		eventID int
		mock    func()
		want    want
	}{
		{
			name:    "#1 OK",
			eventID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"event_id", "topic_id", "name", "year", "date"}).
					AddRow(testEvent.EventID, testEvent.TopicID, testEvent.Name, testEvent.Year, testEvent.Date)

				mock.ExpectQuery(query).WithArgs(testEvent.EventID).WillReturnRows(rows)
			},
			want: want{
				event: testEvent,
				error: false,
			},
		},
		{
			name:    "#2 NOT FOUND",
			eventID: 20,
			mock: func() {
				rows := sqlmock.NewRows([]string{"event_id", "topic_id", "name", "year", "date"}).
					AddRow(testEvent.EventID, testEvent.TopicID, testEvent.Name, testEvent.Year, testEvent.Date)

				mock.ExpectQuery(query).WithArgs(testEvent.EventID).WillReturnRows(rows)
			},
			want: want{
				event: x.Event{},
				error: true,
			},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			event, err := store.GetEvent(test.eventID)
			if (err != nil) != test.want.error {
				t.Errorf("GetEvent() error = %v, want error %v", err, test.want.error)
				return
			}
			if err == nil && !reflect.DeepEqual(event, test.want.event) {
				t.Errorf("Get() = %v, want %v", event, test.want.event)
			}
		})
	}
}

// TestCountEvents tests getting amount of events.
func TestCountEvents(t *testing.T) {
	// New mock database
	db, mock := NewMock()
	store := &EventStore{db}
	defer db.Close()

	// Expected results
	type want struct {
		eventsCount int
		error       bool
	}

	query := regexp.QuoteMeta(countEventsQuery)

	// Declare test cases
	tests := []struct {
		name string
		mock func()
		want want
	}{
		{
			name: "#1 OK",
			mock: func() {
				rows := mock.NewRows([]string{"event_id", "topic_id", "name", "year", "date"})
				for _, event := range testEvents {
					rows = rows.AddRow(event.EventID, event.TopicID, event.Name, event.Year, event.Date)
				}

				amount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(len(testEvents))
				mock.ExpectQuery(query).WillReturnRows(amount)
			},
			want: want{
				eventsCount: len(testEvents),
				error:       false,
			},
		},
		{
			name: "#2 OK (NO ROWS)",
			mock: func() {
				mock.NewRows([]string{"event_id", "topic_id", "name", "year", "date"})

				amount := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
				mock.ExpectQuery(query).WillReturnRows(amount)
			},
			want: want{
				eventsCount: 0,
				error:       false,
			},
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			eventsCount, err := store.CountEvents()
			if (err != nil) != test.want.error {
				t.Errorf("GetEvent() error = %v, want error %v", err, test.want.error)
				return
			}
			if err == nil && !reflect.DeepEqual(eventsCount, test.want.eventsCount) {
				t.Errorf("Get() = %v, want %v", eventsCount, test.want.eventsCount)
			}
		})
	}
}

// TestCreateEvent tests creating a new event.
func TestCreateEvent(t *testing.T) {

}

// TestUpdateEvent tests updating an existing event.
func TestUpdateEvent(t *testing.T) {

}

// TestDeleteEvent tests deleting an existing event.
func TestDeleteEvent(t *testing.T) {

}
