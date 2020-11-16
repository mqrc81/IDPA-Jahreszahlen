package web

/*
 * TODO Header
 */

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * NewHandler creates a new handler, including routes and middleware
 */
func NewHandler(store backend.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)

	h.Get("/", h.Home())

	h.Route("/users", func(r chi.Router) {
		//r.Get("/register", UsersRegister())
		//r.Get("/login", UsersLogin())
		//r.Get("/{username}", UsersProfile())
		//r.Get("/{username}/edit", UsersEdit())
		//r.Get("/{username}/scoreboard", UsersScoreboard())
		//
		//r.Post("/store", UsersStore())
	})

	h.Route("/topics", func(r chi.Router) {
		r.Get("/", h.TopicsList())
		r.Get("/new", h.TopicsCreate())
		r.Post("/store", h.TopicsStore())
		r.Post("/{topicID}/delete", h.TopicsDelete())
		r.Get("/{topicID}/edit", h.TopicsEdit())
		r.Get("/{topicID}", h.TopicsShow())
		r.Get("/{topicID}/play", h.TopicsPlay())
		//r.Get("/{topicID}/scoreboard", TopicsScoreboard())

		r.Get("/{topicID}/events/new", h.EventsCreate())
		r.Post("/{topicID}/events/store", h.EventsStore())
		r.Post("/{topicID}/events/{eventID}/delete", h.EventsDelete())
	})

	return h
}

/*
 * Handler consists of the chi-multiplexer and a store with functions for topics and events
 */
type Handler struct {
	*chi.Mux
	store backend.Store
}

/*
 * Home TODO
 */
func (h *Handler) Home() http.HandlerFunc {
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(``))

	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

/*
 * TopicsList is a GET method that lists all topics
 */
func (h *Handler) TopicsList() http.HandlerFunc {
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
 * TopicsCreate is a GET method for a form to create a new topic
 */
func (h *Handler) TopicsCreate() http.HandlerFunc {
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
 * TopicsStore is a POST method that stores topic created
 */
func (h *Handler) TopicsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve variables from form (TopicsCreate)
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
 * TopicsDelete is a POST method that deletes a topic
 */
func (h *Handler) TopicsDelete() http.HandlerFunc {
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
 * TopicsShow is a GET method that shows a specific topic with options to play, see leaderboard, (edit if admin)
 */
func (h *Handler) TopicsShow() http.HandlerFunc {
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
 * TopicsPlay is a GET method that goes through the 3 phases of the quiz and stores the user's score TODO
 */
func (h *Handler) TopicsPlay() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Topic       backend.Topic
		Events      []backend.Event
		EventsCount int
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(`TODO`)) // TODO

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
 * TopicsEdit is a GET method with the option to edit a specific topic and its events
 */
func (h *Handler) TopicsEdit() http.HandlerFunc {
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
 * EventsCreate is a GET method for a form to create a new event
 */
func (h *Handler) EventsCreate() http.HandlerFunc {
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(eventsCreateHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * EventsStore is a POST method that stores event created
 */
func (h *Handler) EventsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Retrieve variables from form (EventsCreate)
		title := r.FormValue("title")
		year, _ := strconv.Atoi(r.FormValue("year"))

		// Execute SQL statement
		if err := h.store.CreateEvent(&backend.Event{
			EventID: 0,
			TopicID: topicID,
			Title:   title,
			Year:    year,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

		// Redirect to list of topics
		http.Redirect(w, r, "/topics", http.StatusFound)
	}
}

/*
 * EventsDelete is a POST method that deletes an event
 */
func (h *Handler) EventsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve event ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Retrieve event ID from URL
		eventID, _ := strconv.Atoi(chi.URLParam(r, "eventID"))

		// Execute SQL statement
		if err := h.store.DeleteEvent(eventID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(w, r, "/topics/" + strconv.Itoa(topicID) + "/edit", http.StatusFound)
	}
}
