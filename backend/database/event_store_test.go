// Collection of tests for the database access layer of functions evolving
// around events.

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
	store := &EventStore{DB: db}
	defer db.Close()

	query := regexp.QuoteMeta(getEventQuery)

	// Declare test cases
	tests := []struct {
		name      string
		eventID   int
		mock      func()
		wantEvent x.Event
		wantError bool
	}{
		{
			name:    "#1 OK",
			eventID: 1,
			mock: func() {
				rows := sqlmock.NewRows([]string{"event_id", "topic_id", "name", "year", "date"}).
					AddRow(testEvent.EventID, testEvent.TopicID, testEvent.Name, testEvent.Year, testEvent.Date)

				mock.ExpectQuery(query).WithArgs(testEvent.EventID).WillReturnRows(rows)
			},
			wantEvent: testEvent,
			wantError: false,
		},
		{
			name:    "#2 NOT FOUND",
			eventID: 20,
			mock: func() {
				rows := sqlmock.NewRows([]string{"event_id", "topic_id", "name", "year", "date"}).
					AddRow(testEvent.EventID, testEvent.TopicID, testEvent.Name, testEvent.Year, testEvent.Date)

				mock.ExpectQuery(query).WithArgs(testEvent.EventID).WillReturnRows(rows)
			},
			wantEvent: x.Event{},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			event, err := store.GetEvent(test.eventID)
			if (err != nil) != test.wantError {
				t.Errorf("GetEvent() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(event, test.wantEvent) {
				t.Errorf("GetEvent() = %v, want %v", event, test.wantEvent)
			}
		})
	}
}

// TestCountEvents tests getting amount of events.
func TestCountEvents(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &EventStore{DB: db}
	defer db.Close()

	query := regexp.QuoteMeta(countEventsQuery)

	// Declare test cases
	tests := []struct {
		name            string
		mock            func()
		wantEventsCount int
		wantError       bool
	}{
		{
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(len(testEvents))

				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			wantEventsCount: len(testEvents),
			wantError:       false,
		},
		{
			name: "#2 OK (NO ROWS)",
			mock: func() {
				rows := sqlmock.NewRows([]string{"COUNT(*)"}).AddRow(0)
				mock.ExpectQuery(query).WillReturnRows(rows)
			},
			wantEventsCount: 0,
			wantError:       false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			eventsCount, err := store.CountEvents()
			if (err != nil) != test.wantError {
				t.Errorf("CountEvents() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(eventsCount, test.wantEventsCount) {
				t.Errorf("CountEvents() = %v, want %v", eventsCount, test.wantEventsCount)
			}
		})
	}
}

// TestCreateEvent tests creating a new event.
func TestCreateEvent(t *testing.T) {

	// New mock database
	db, _ := NewMock()
	store := EventStore{DB: db}
	defer db.Close()

	_ = regexp.QuoteMeta(createEventQuery)

	// Declare test cases
	tests := []struct {
		name      string
		event     x.Event
		mock      func()
		wantError bool
	}{
		{
			name:  "#1 OK",
			event: testEvent,
			mock: func() {
				// TODO
			},
			wantError: false,
		},
		// TODO
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err := store.CreateEvent(&test.event)
			if (err != nil) != test.wantError {
				t.Errorf("CreateEvent() error = %v, want error %v", err, test.wantError)
				return
			}
		})
	}
}

// TestUpdateEvent tests updating an existing event.
func TestUpdateEvent(t *testing.T) {

	// New mock database
	db, _ := NewMock()
	store := EventStore{DB: db}
	defer db.Close()

	_ = regexp.QuoteMeta(createEventQuery)

	// Declare test cases
	tests := []struct {
		name      string
		event     x.Event
		mock      func()
		wantError bool
	}{
		{
			name:  "#1 OK",
			event: testEvent,
			mock: func() {
				// TODO
			},
			wantError: false,
		},
		// TODO
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err := store.UpdateEvent(&test.event)
			if (err != nil) != test.wantError {
				t.Errorf("UpdateEvent() error = %v, want error %v", err, test.wantError)
				return
			}
		})
	}
}

// TestDeleteEvent tests deleting an existing event.
func TestDeleteEvent(t *testing.T) {

	// New mock database
	db, _ := NewMock()
	store := EventStore{DB: db}
	defer db.Close()

	_ = regexp.QuoteMeta(deleteEventQuery)

	// Declare test cases
	tests := []struct {
		name      string
		eventID   int
		mock      func()
		wantError bool
	}{
		{
			name:    "#1 OK",
			eventID: testEvent.EventID,
			mock: func() {
				// TODO
			},
			wantError: false,
		},
		// TODO
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.mock()
			err := store.DeleteEvent(test.eventID)
			if (err != nil) != test.wantError {
				t.Errorf("UpdateEvent() error = %v, want error %v", err, test.wantError)
				return
			}
		})
	}
}
