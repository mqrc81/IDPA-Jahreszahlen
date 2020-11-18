package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type TopicHandler struct {
	store backend.Store
	sessions *scs.SessionManager
}

/*
 * List is a GET method that lists all topics
 */
func (h *TopicHandler) List() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Topics []backend.Topic
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(topicsListHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Execute SQL statement and return slice of topics
		tt, err := h.store.Topics()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-template
		if err := tmpl.Execute(w, data{Topics: tt}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Create is a GET method for a form to create a new topic
 */
func (h *TopicHandler) Create() http.HandlerFunc {
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(topicsCreateHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Store is a POST method that stores topic created
 */
func (h *TopicHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve variables from form (Create)
		title := r.FormValue("title")
		startYear, _ := strconv.Atoi(r.FormValue("start_year"))
		endYear, _ := strconv.Atoi(r.FormValue("end_year"))
		description := r.FormValue("description")

		// Execute SQL statement
		if err := h.store.CreateTopic(&backend.Topic{
			TopicID:     0,
			Title:       title,
			StartYear:   startYear,
			EndYear:     endYear,
			Description: description,
			PlayCount:   0,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(w, r, "/topics", http.StatusFound)
	}
}

/*
 * Delete is a POST method that deletes a topic
 */
func (h *TopicHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve TopicID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Execute SQL statement
		if err := h.store.DeleteTopic(topicID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(w, r, "/topics", http.StatusFound)
	}
}

/*
 * Edit is a GET method with the option to edit a specific topic and its events
 */
func (h *TopicHandler) Edit() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Topic  backend.Topic
		Events []backend.Event
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(topicsEditHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Execute SQL statement and return topic
		t, err := h.store.Topic(topicID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement and return events
		ee, err := h.store.EventsByTopic(topicID, false)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-template
		if err := tmpl.Execute(w, data{Topic: t, Events: ee}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Show is a GET method that shows a specific topic with options to play, see leaderboard, (edit if admin)
 */
func (h *TopicHandler) Show() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Topic backend.Topic
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(topicsShowHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve TopicID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Execute SQL statement and return topic
		t, err := h.store.Topic(topicID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-template
		if err := tmpl.Execute(w, data{Topic: t}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Play is a GET method that goes through the 3 phases of the quiz and stores the user's score TODO
 */
func (h *TopicHandler) Play() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Topic       backend.Topic
		Events      []backend.Event
		EventsCount int
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(`<h1>TODO</h1>`)) // TODO

	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Execute SQL statement and return topic
		t, err := h.store.Topic(topicID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement and return events
		ee, err := h.store.EventsByTopic(topicID, true)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-template
		if err := tmpl.Execute(w, data{Topic: t, Events: ee, EventsCount: len(ee)}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Scoreboard is a GET method that lists scores of a topic sorted by points
 */
func (h *TopicHandler) Scoreboard() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Scores []backend.Score
		TopicName string
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Funcs(funcMap).Parse(topicsScoreboardHTML))
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Execute SQL statement and return slice of scores
		ss, err := h.store.ScoresByTopic(topicID, 50)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement and return topic
		t, err := h.store.Topic(topicID)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Execute HTML-template
		if err := tmpl.Execute(w, data{Scores: ss, TopicName: t.Title}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}