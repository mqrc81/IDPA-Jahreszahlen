package web

// event_handler.go
// Contains all HTTP-handlers for pages evolving around events.

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// EventHandler
// Object for handlers to access sessions and database.
type EventHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// List
// A GET-method that any admin can call. It lists all scores, ranked by points,
// with the ability to filter scores by topic and/or user.
func (h *EventHandler) List() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topic  backend.Topic
		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/events_list.html"))
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute SQL statement to get a topic
		topic, err := h.store.Topic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get events
		events, err := h.store.EventsByTopic(topicID, false)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Topic:       topic,
			Events:      events,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Create
// A GET-method that any admin can call. It renders a form, in which values
// for a new event can be entered.
func (h *EventHandler) Create() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topic backend.Topic
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/events_create.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute SQL statement to get a topic
		topic, err := h.store.Topic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Topic:       topic,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Store
// A POST-method. It validates the form from Create and redirects to Create in
// case of an invalid input with corresponding error message. In case of valid
// form, it stores the new event in the database and redirects to the edit-page
// of the event's topic.
func (h *EventHandler) Store() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve variables from form (Create)
		year, _ := strconv.Atoi(req.FormValue("year"))
		form := EventForm{
			Title: req.FormValue("title"),
			Year:  year,
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to create an event
		if err := h.store.CreateEvent(&backend.Event{
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
		http.Redirect(res, req, "/topics/"+topicIDstr+"/events", http.StatusFound)
	}
}

// Delete
// A POST-method. It deletes a certain event and redirects to edit-page of the
// event's topic.
func (h *EventHandler) Delete() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve event ID from URL
		topicID := chi.URLParam(req, "topicID")

		// Retrieve event ID from URL
		eventID, _ := strconv.Atoi(chi.URLParam(req, "eventID"))

		// Execute SQL statement to delete an event
		if err := h.store.DeleteEvent(eventID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicID+"/edit", http.StatusFound)
	}
}
