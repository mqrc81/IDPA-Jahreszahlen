// Collection of tests for the database access layer of functions evolving
// around events.

package database

import (
	_ "database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/util"
)

var (
	// tEvent is a mock event for testing purposes
	tEvent = x.Event{
		EventID: 1,
		TopicID: 1,
		Name:    "Test Event 1",
		Year:    1800,
		Date:    util.Date(1800, 1, 1),
	}
)

// TestGetEvent tests getting an event by ID.
func TestGetEvent(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &EventStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM events"

	table := []string{"event_id", "topic_id", "name", "year", "date"}

	// Declare test cases
	tests := []struct {
		name      string
		eventID   int
		mock      func(eventID int)
		wantEvent x.Event
		wantError bool
	}{
		{
			// When everything works as intended
			name:    "#1 OK",
			eventID: 1,
			mock: func(eventID int) {
				rows := sqlmock.NewRows(table).
					AddRow(tEvent.EventID, tEvent.TopicID, tEvent.Name, tEvent.Year, tEvent.Date)

				mock.ExpectQuery(queryMatch).WithArgs(eventID).WillReturnRows(rows)
			},
			wantEvent: tEvent,
			wantError: false,
		},
		{
			// When event with given event ID doesn't exist
			name:    "#2 NOT FOUND",
			eventID: 0,
			mock: func(eventID int) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(eventID).WillReturnRows(rows)
			},
			wantEvent: x.Event{},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.eventID)

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

	queryMatch := "SELECT COUNT((.+)) FROM events"

	table := []string{"COUNT(*)"}

	// Declare test cases
	tests := []struct {
		name            string
		mock            func()
		wantEventsCount int
		wantError       bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows(table).AddRow(3)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantEventsCount: 3,
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
	db, mock := NewMock()
	store := &EventStore{DB: db}
	defer db.Close()

	queryMatch := "INSERT INTO events"

	// Declare test cases
	tests := []struct {
		name      string
		event     x.Event
		mock      func(event x.Event)
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			event: tEvent,
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.TopicID, event.Name, event.Year, event.Date).
					WillReturnResult(sqlmock.NewResult(int64(event.TopicID), 1))
			},
			wantError: false,
		},
		{
			// When topic with given topic ID doesn't exist
			name: "#2 TOPIC NOT FOUND",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: 0,
				Name:    tEvent.Name,
				Year:    tEvent.Year,
				Date:    tEvent.Date,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.TopicID, event.Name, event.Year, event.Date).
					WillReturnError(errors.New("topic does not exist"))
			},
			wantError: true,
		},
		{
			// When title is missing
			name: "#3 NAME MISSING",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: tEvent.TopicID,
				Year:    tEvent.Year,
				Date:    tEvent.Date,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.TopicID, event.Name, event.Year, event.Date).
					WillReturnError(errors.New("name can not be empty"))
			},
			wantError: true,
		},
		{
			// When year is missing
			name: "#4 YEAR MISSING",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: tEvent.TopicID,
				Name:    tEvent.Name,
				Date:    tEvent.Date,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.TopicID, event.Name, event.Year, event.Date).
					WillReturnError(errors.New("year can not be empty"))
			},
			wantError: true,
		},
		{
			// When date is missing
			name: "#5 DATE MISSING",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: tEvent.TopicID,
				Name:    tEvent.Name,
				Year:    tEvent.Year,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.TopicID, event.Name, event.Year, event.Date).
					WillReturnError(errors.New("date can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.event)

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
	db, mock := NewMock()
	store := &EventStore{DB: db}
	defer db.Close()

	queryMatch := "UPDATE events"

	// Declare test cases
	tests := []struct {
		name      string
		event     x.Event
		mock      func(event x.Event)
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			event: tEvent,
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(tEvent.Name, tEvent.Year, tEvent.Date, tEvent.EventID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			// When event with given event ID doesn't exist
			name: "#2 NOT FOUND",
			event: x.Event{
				EventID: 0,
				TopicID: tEvent.TopicID,
				Name:    tEvent.Name,
				Year:    tEvent.Year,
				Date:    tEvent.Date,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(tEvent.Name, tEvent.Year, tEvent.Date, tEvent.EventID).
					WillReturnError(errors.New("event with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When title is missing
			name: "#3 NAME MISSING",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: tEvent.TopicID,
				Year:    tEvent.Year,
				Date:    tEvent.Date,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.Name, event.Year, event.Date, event.EventID).
					WillReturnError(errors.New("name can not be empty"))
			},
			wantError: true,
		},
		{
			// When year is missing
			name: "#4 YEAR MISSING",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: tEvent.TopicID,
				Name:    tEvent.Name,
				Date:    tEvent.Date,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.Name, event.Year, event.Date, event.EventID).
					WillReturnError(errors.New("year can not be empty"))
			},
			wantError: true,
		},
		{
			// When date is missing
			name: "#5 DATE MISSING",
			event: x.Event{
				EventID: tEvent.EventID,
				TopicID: tEvent.TopicID,
				Name:    tEvent.Name,
				Year:    tEvent.Year,
			},
			mock: func(event x.Event) {
				mock.ExpectExec(queryMatch).WithArgs(event.Name, event.Year, event.Date, event.EventID).
					WillReturnError(errors.New("date can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.event)

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
	db, mock := NewMock()
	store := &EventStore{DB: db}
	defer db.Close()

	queryMatch := "DELETE FROM events"

	// Declare test cases
	tests := []struct {
		name      string
		eventID   int
		mock      func(eventID int)
		wantError bool
	}{
		{
			// When everything works as intended
			name:    "#1 OK",
			eventID: tEvent.EventID,
			mock: func(eventID int) {
				mock.ExpectExec(queryMatch).WithArgs(eventID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			// When event with given event ID doesn't exist
			name:    "#2 NOT FOUND",
			eventID: 0,
			mock: func(eventID int) {
				mock.ExpectExec(queryMatch).WithArgs(eventID).
					WillReturnError(errors.New("event with given id does not exist"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.eventID)

			err := store.DeleteEvent(test.eventID)

			if (err != nil) != test.wantError {
				t.Errorf("DeleteEvent() error = %v, want error %v", err, test.wantError)
				return
			}
		})
	}
}
