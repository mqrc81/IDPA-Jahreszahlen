package web

/*
 * event_handler.go contains HTTP-handler functions for events
 */

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * EventHandler handles sessions, CSRF-protection and database access for events
 */
type EventHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

/*
 * Create is a GET method for a form to create a new event
 */
func (h *EventHandler) Create() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData

		TopicID int
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(eventsCreateHTML))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			TopicID:     topicID,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Store is a POST method that stores event created
 */
func (h *EventHandler) Store() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve variables from form (Create)
		year, _ := strconv.Atoi(req.FormValue("year"))
		form := CreateEventForm{
			Title: req.FormValue("title"),
			Year:  year,
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement
		if err := h.store.CreateEvent(&backend.Event{
			EventID: 0,
			TopicID: topicID,
			Title:   form.Title,
			Year:    form.Year,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Adds flash message
		h.sessions.Put(req.Context(), "flash", "Ereignis wurde erfolgreich erstellt.")

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicIDstr+"/edit", http.StatusFound)
	}
}

/*
 * Delete is a POST method that deletes an event
 */
func (h *EventHandler) Delete() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve event ID from URL
		topicID := chi.URLParam(req, "topicID")

		// Retrieve event ID from URL
		eventID, _ := strconv.Atoi(chi.URLParam(req, "eventID"))

		// Execute SQL statement
		if err := h.store.DeleteEvent(eventID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicID+"/edit", http.StatusFound)
	}
}
