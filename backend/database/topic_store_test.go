// Collection of tests for the database access layer of functions evolving
// around topics.

package database

import (
	"reflect"
	"testing"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend/util"
)

var (
	// tTopic is a mock topic for testing purposes
	tTopic = x.Topic{
		TopicID:     1,
		Name:        "Test Topic 1",
		StartYear:   1800,
		EndYear:     1900,
		Description: "Test Description",
		Events: []x.Event{
			tEvent,
			{
				EventID: 2,
				TopicID: 1,
				Name:    "Test Event 2",
				Year:    1850,
				Date:    util.Date(1850, 1, 1),
			},
		},
		ScoresCount: 15,
		EventsCount: 2,
	}

	// tTopics is a mock array of topics for testing purposes
	tTopics = []x.Topic{
		tTopic,
		{
			TopicID:     2,
			Name:        "Test Topic 2",
			StartYear:   1700,
			EndYear:     1800,
			Description: "Test Description",
			Events: []x.Event{
				{
					EventID: 3,
					TopicID: 2,
					Name:    "Test Event 3",
					Year:    1750,
					Date:    util.Date(1750, 1, 1),
				},
			},
			ScoresCount: 30,
			EventsCount: 1,
		},
	}
)

// TestGetTopic tests getting a topic by ID.
func TestGetTopic(t *testing.T) {

	// New mock database
	db, mock := NewMock()
	store := &TopicStore{DB: db}
	defer db.Close()

	queryMatch := "SELECT (.+) FROM topics"
	queryMatchEvents := "SELECT (.+) FROM events"

	table := []string{"topic_id", "name", "start_year", "end_year", "description", "scores_count", "events_count"}
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
				rows := mock.NewRows(table).
					AddRow(tTopic.TopicID, tTopic.Name, tTopic.StartYear, tTopic.EndYear, tTopic.Description,
						tTopic.ScoresCount, tTopic.EventsCount)

				mock.ExpectQuery(queryMatch).WithArgs(topicID).WillReturnRows(rows)

				rowsEvents := mock.NewRows(tableEvents)
				for _, e := range tTopic.Events {
					rowsEvents = rowsEvents.AddRow(e.EventID, e.TopicID, e.Name, e.Year, e.Date)
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
				rows := mock.NewRows(table)

				mock.ExpectQuery(queryMatch).WithArgs(topicID).WillReturnRows(rows)

				rowsEvents := mock.NewRows(tableEvents)
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

}

// TestCountTopics tests getting amount of topics.
func TestCountTopics(t *testing.T) {

}

// TestCreateTopic tests creating a new topic.
func TestCreateTopic(t *testing.T) {

}

// TestUpdateTopic tests updating an existing topic.
func TestUpdateTopic(t *testing.T) {

}

// TestDeleteTopic tests deleting an existing topic.
func TestDeleteTopic(t *testing.T) {

}
