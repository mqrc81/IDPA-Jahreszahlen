package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

type EventHandler struct {
	store backend.Store
	sessions *scs.SessionManager
}

/*
 * Create is a GET method for a form to create a new event
 */
func (h *EventHandler) Create() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		TopicID int
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(eventsCreateHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(r, "topicID"))

		// Execute HTML-template
		if err := tmpl.Execute(w, data{topicID}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Store is a POST method that stores event created
 */
func (h *EventHandler) Store() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve topic ID from URL
		topicIDstr := chi.URLParam(r, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve variables from form (Create)
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
			return
		}

		// Redirect to list of topics
		http.Redirect(w, r, "/topics/" +topicIDstr + "/edit", http.StatusFound)
	}
}

/*
 * Delete is a POST method that deletes an event
 */
func (h *EventHandler) Delete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Retrieve event ID from URL
		topicID := chi.URLParam(r, "topicID")

		// Retrieve event ID from URL
		eventID, _ := strconv.Atoi(chi.URLParam(r, "eventID"))

		// Execute SQL statement
		if err := h.store.DeleteEvent(eventID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(w, r, "/topics/" + topicID + "/edit", http.StatusFound)
	}
}

