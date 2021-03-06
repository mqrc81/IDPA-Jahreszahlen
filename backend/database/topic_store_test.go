// Collection of tests for the database access layer of functions evolving
// around topics.

package database

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// tTopic is a mock topic for testing purposes
	tTopic = x.Topic{
		TopicID:     1,
		Name:        "Test Topic 1",
		StartYear:   1800,
		EndYear:     1900,
		Description: "Test Description",
		Image:       "https://test-image.png",
		Events: []x.Event{
			tEvent,
			{
				EventID: 2,
				TopicID: 1,
				Name:    "Test Event 2",
				Year:    1850,
				Date:    time.Date(1850, 1, 1, 0, 0, 0, 0, time.UTC),
			},
		},
		ScoresCount: 15,
		EventsCount: 2,
	}

	// tTopics is a mock array of topics for testing purposes
	tTopics = []x.Topic{
		{
			TopicID:     tTopic.TopicID,
			Name:        tTopic.Name,
			StartYear:   tTopic.StartYear,
			EndYear:     tTopic.EndYear,
			Description: tTopic.Description,
			Image:       tTopic.Image,
			ScoresCount: tTopic.ScoresCount,
			EventsCount: tTopic.EventsCount,
		},
		{
			TopicID:     2,
			Name:        "Test Topic 2",
			StartYear:   1700,
			EndYear:     1800,
			Description: "Test Description",
			Image:       "https://test-image-2.png",
			ScoresCount: 30,
			EventsCount: 15,
		},
	}

	// nilTopics is a nil slice of topics, since "var t []Topic" is a nil slice
	// but "t := []Topic{}" is an empty slice (so we can't use the latter for
	// this use case)
	nilTopics []x.Topic
)

// TestGetTopic tests getting a topic by ID.
func TestGetTopic(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TopicStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM topics"
	queryMatchEvents := "SELECT (.+) FROM events"

	table := []string{"topic_id", "name", "start_year", "end_year", "description", "image", "scores_count",
		"events_count"}
	tableEvents := []string{"event_id", "topic_id", "name", "year", "date"}

	// Declare test cases
	tests := []struct {
		name      string
		topicID   int
		mock      func(topicID int)
		wantTopic x.Topic
		wantError bool
	}{
		{
			// When everything works as expected
			name:    "#1 OK",
			topicID: tTopic.TopicID,
			mock: func(topicID int) {
				rows := sqlmock.NewRows(table).
					AddRow(tTopic.TopicID, tTopic.Name, tTopic.StartYear, tTopic.EndYear, tTopic.Description,
						tTopic.Image, tTopic.ScoresCount, tTopic.EventsCount)

				mock.ExpectQuery(queryMatch).WithArgs(topicID).WillReturnRows(rows)

				rowsEvents := sqlmock.NewRows(tableEvents)
				for _, event := range tTopic.Events {
					rowsEvents = rowsEvents.AddRow(event.EventID, event.TopicID, event.Name, event.Year, event.Date)
				}
				mock.ExpectQuery(queryMatchEvents).WithArgs(topicID).WillReturnRows(rowsEvents)
			},
			wantTopic: tTopic,
			wantError: false,
		},
		{
			// When topic with given topic ID doesn't exist
			name:    "#2 NOT FOUND",
			topicID: 0,
			mock: func(topicID int) {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(topicID).WillReturnRows(rows)

				rowsEvents := sqlmock.NewRows(tableEvents)
				mock.ExpectQuery(queryMatchEvents).WithArgs(topicID).WillReturnRows(rowsEvents)
			},
			wantTopic: x.Topic{},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.topicID)

			topic, err := store.GetTopic(test.topicID)

			if (err != nil) != test.wantError {
				t.Errorf("GetTopic() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(topic, test.wantTopic) {
				t.Errorf("GetTopic() = %v, want %v", topic, test.wantTopic)
			}
		})
	}
}

// TestGetTopics tests getting all topics.
func TestGetTopics(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TopicStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM topics"

	table := []string{"topic_id", "name", "start_year", "end_year", "description", "image", "scores_count",
		"events_count"}

	// Declare test cases
	tests := []struct {
		name       string
		mock       func()
		wantTopics []x.Topic
		wantError  bool
	}{
		{
			// When everything works as intended
			name: "#1 OK",
			mock: func() {
				rows := sqlmock.NewRows(table)
				for _, topic := range tTopics {
					rows = rows.AddRow(topic.TopicID, topic.Name, topic.StartYear, topic.EndYear, topic.Description,
						topic.Image, topic.ScoresCount, topic.EventsCount)

					mock.ExpectQuery(queryMatch).WillReturnRows(rows)
				}
			},
			wantTopics: tTopics,
			wantError:  false,
		},
		{
			// When topics table is empty
			name: "#2 OK (NO ROWS)",
			mock: func() {
				rows := sqlmock.NewRows(table)

				mock.ExpectQuery(queryMatch).WillReturnRows(rows)
			},
			wantTopics: nilTopics,
			wantError:  false,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock()

			topics, err := store.GetTopics()

			if (err != nil) != test.wantError {
				t.Errorf("GetTopics() error = %v, want error %v", err, test.wantError)
				return
			}
			if err == nil && !reflect.DeepEqual(topics, test.wantTopics) {
				t.Errorf("GetTopics() = %v, want %v", topics, test.wantTopics)
			}
		})
	}
}

// TestCreateTopic tests creating a new topic.
func TestCreateTopic(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TopicStore{DB: db}
	defer db.Close()

	queryMatch := "INSERT INTO topics"

	// Declare test cases
	tests := []struct {
		name      string
		topic     x.Topic
		mock      func(topic x.Topic)
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			topic: tTopic,
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					topic.Image).
					WillReturnResult(sqlmock.NewResult(int64(topic.TopicID), 1))
			},
			wantError: false,
		},
		{
			// When name is missing
			name: "#2 NAME MISSING",
			topic: x.Topic{
				StartYear:   tTopic.StartYear,
				EndYear:     tTopic.EndYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					topic.Image).
					WillReturnError(errors.New("name can not be empty"))
			},
			wantError: true,
		},
		{
			// When start-year is missing
			name: "#3 START-YEAR MISSING",
			topic: x.Topic{
				Name:        tTopic.Name,
				EndYear:     tTopic.EndYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					topic.Image).
					WillReturnError(errors.New("start-year can not be empty"))
			},
			wantError: true,
		},
		{
			// When end-year is missing
			name: "#4 END-YEAR MISSING",
			topic: x.Topic{
				Name:        tTopic.Name,
				StartYear:   tTopic.StartYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					topic.Image).
					WillReturnError(errors.New("end-year can not be empty"))
			},
			wantError: true,
		},
		{
			// When description is missing
			name: "#5 OK (DESCRIPTION MISSING)",
			topic: x.Topic{
				TopicID:   tTopic.TopicID,
				Name:      tTopic.Name,
				StartYear: tTopic.StartYear,
				EndYear:   tTopic.EndYear,
				Image:     tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					topic.Image).
					WillReturnResult(sqlmock.NewResult(int64(topic.TopicID), 1))
			},
			wantError: false,
		},
		{
			// When image is missing
			name: "#5 IMAGE MISSING",
			topic: x.Topic{
				TopicID:     tTopic.TopicID,
				Name:        tTopic.Name,
				StartYear:   tTopic.StartYear,
				EndYear:     tTopic.EndYear,
				Description: tTopic.Description,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					topic.Image).
					WillReturnError(errors.New("image can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.topic)

			err := store.CreateTopic(&test.topic)

			if (err != nil) != test.wantError {
				t.Errorf("CreateTopic() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}

// TestUpdateTopic tests updating an existing topic.
func TestUpdateTopic(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TopicStore{DB: db}
	defer db.Close()

	queryMatch := "UPDATE topics"

	// Declare test cases
	tests := []struct {
		name      string
		topic     x.Topic
		mock      func(topic x.Topic)
		wantError bool
	}{
		{
			// When everything works as intended
			name:  "#1 OK",
			topic: tTopic,
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnResult(sqlmock.NewResult(int64(topic.TopicID), 1))
			},
			wantError: false,
		},
		{
			// When topic with given topic ID doesn't exist
			name: "#2 NOT FOUND",
			topic: x.Topic{
				TopicID:     0,
				Name:        tTopic.Name,
				StartYear:   tTopic.StartYear,
				EndYear:     tTopic.EndYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnError(errors.New("topic with given id does not exist"))
			},
			wantError: true,
		},
		{
			// When name is missing
			name: "#3 NAME MISSING",
			topic: x.Topic{
				TopicID:     tTopic.TopicID,
				StartYear:   tTopic.StartYear,
				EndYear:     tTopic.EndYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnError(errors.New("name can not be empty"))
			},
			wantError: true,
		},
		{
			// When start-year is missing
			name: "#4 START-YEAR MISSING",
			topic: x.Topic{
				TopicID:     tTopic.TopicID,
				Name:        tTopic.Name,
				EndYear:     tTopic.EndYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnError(errors.New("start-year can not be empty"))
			},
			wantError: true,
		},
		{
			// When end-year is missing
			name: "#5 END-YEAR MISSING",
			topic: x.Topic{
				TopicID:     tTopic.TopicID,
				Name:        tTopic.Name,
				StartYear:   tTopic.StartYear,
				Description: tTopic.Description,
				Image:       tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnError(errors.New("end-year can not be empty"))
			},
			wantError: true,
		},
		{
			// When description is missing
			name: "#6 OK (DESCRIPTION MISSING)",
			topic: x.Topic{
				TopicID:   tTopic.TopicID,
				Name:      tTopic.Name,
				StartYear: tTopic.StartYear,
				EndYear:   tTopic.EndYear,
				Image:     tTopic.Image,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnResult(sqlmock.NewResult(int64(topic.TopicID), 1))
			},
			wantError: false,
		},
		{
			// When image is missing
			name: "#5 IMAGE MISSING",
			topic: x.Topic{
				TopicID:     tTopic.TopicID,
				Name:        tTopic.Name,
				StartYear:   tTopic.StartYear,
				Description: tTopic.Description,
			},
			mock: func(topic x.Topic) {
				mock.ExpectExec(queryMatch).WithArgs(topic.Name, topic.StartYear, topic.EndYear, topic.Description,
					tTopic.Image, topic.TopicID).
					WillReturnError(errors.New("image can not be empty"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.topic)

			err := store.UpdateTopic(&test.topic)

			if (err != nil) != test.wantError {
				t.Errorf("UpdateTopic() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}

// TestDeleteTopic tests deleting an existing topic.
func TestDeleteTopic(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TopicStore{DB: db}
	defer db.Close()

	queryMatch := "DELETE FROM topics"

	// Declare test cases
	tests := []struct {
		name      string
		topicID   int
		mock      func(topicID int)
		wantError bool
	}{
		{
			// When everything works as intended
			name:    "#1 OK",
			topicID: tTopic.TopicID,
			mock: func(topicID int) {
				mock.ExpectExec(queryMatch).WithArgs(topicID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			wantError: false,
		},
		{
			// When topic with given topic ID doesn't exist
			name:    "#2 NOT FOUND",
			topicID: 0,
			mock: func(topicID int) {
				mock.ExpectExec(queryMatch).WithArgs(topicID).
					WillReturnError(errors.New("topic with given id does not exist"))
			},
			wantError: true,
		},
	}

	// Run tests
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			test.mock(test.topicID)

			err := store.DeleteTopic(test.topicID)

			if (err != nil) != test.wantError {
				t.Errorf("DeleteTopic() error = %v, want error %v", err, test.wantError)
			}
		})
	}
}
