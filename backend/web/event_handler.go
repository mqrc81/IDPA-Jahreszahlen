package web

/*
 * Contains all HTTP-handler functions for pages evolving around events.
 */

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// EventHandler is the object for handlers to access sessions and database.
type EventHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// List is a GET-method that any user can call. It lists all scores, ranked by
// points, with the ability to filter scores by topic and/or user.
func (handler *EventHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topic  backend.Topic
		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/events_list.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Benutzer eingeloggt sein, um all Ereignisse eines Themas aufzulisten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Retrieve topic ID from URL parameters
		topicID, err := strconv.Atoi(chi.URLParam(req, "topicID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Execute SQL statement to get a topic
		topic, err := handler.store.GetTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topic:       topic,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Create is a GET-method that any admin can call. It displays a form, in which
// values for a new event can be entered.
func (handler *EventHandler) Create() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topic backend.Topic
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/events_create.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(backend.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Admin eingeloggt sein, um ein neues Ereignis zu erstellen.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Retrieve topic ID from URL parameters
		topicID, err := strconv.Atoi(chi.URLParam(req, "topicID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Execute SQL statement to get a topic
		topic, err := handler.store.GetTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topic:       topic,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CreateStore is a POST-method. It validates the form from Create and
// redirects to Create in case of an invalid input with corresponding error
// message. In case of valid form, it stores the new event in the database and
// redirects to the edit-page of the event's topic.
func (handler *EventHandler) CreateStore() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve variables from form (Create)
		year, _ := strconv.Atoi(req.FormValue("year"))
		form := EventForm{
			Name: req.FormValue("name"),
			Year: year,
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to create an event
		if err := handler.store.CreateEvent(&backend.Event{
			TopicID: topicID,
			Name:    form.Name,
			Year:    form.Year,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Adds flash message
		handler.sessions.Put(req.Context(), "flash_success", "Ereignis wurde erfolgreich erstellt.")

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicIDstr+"/events", http.StatusFound)
	}
}

// Delete is a POST-method. It deletes an event and redirects to the list of
// events.
func (handler *EventHandler) Delete() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve event ID from URL parameters
		topicID := chi.URLParam(req, "topicID")

		// Retrieve event ID from URL parameters
		eventID, _ := strconv.Atoi(chi.URLParam(req, "eventID"))

		// Execute SQL statement to delete an event
		if err := handler.store.DeleteEvent(eventID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicID+"/events", http.StatusFound)
	}
}

// Edit is a GET-method that any admin can call. It displays a form in which
// values for updating the current event can be entered.
func (handler *EventHandler) Edit() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Event backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/layout.html",
		"frontend/pages/events_edit.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(backend.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Admin eingeloggt sein, um ein Ereignis zu bearbeiten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Retrieve event ID from URL parameters
		eventID, err := strconv.Atoi(chi.URLParam(req, "eventID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Execute SQL statement to get topic
		event, err := handler.store.GetEvent(eventID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Event:       event,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditStore is a POST-method. It validates the form from Edit and redirects to
// Edit in case of an invalid input with corresponding error message. In case
// of valid form, it stores the topic in the database and redirects to List.
func (handler *EventHandler) EditStore() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve values from form (Edit)
		year, _ := strconv.Atoi(req.FormValue("year"))
		form := EventForm{
			Name: req.FormValue("name"),
			Year: year,
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to update event
		if err := handler.store.UpdateEvent(&backend.Event{
			TopicID: topicID,
			Name:    form.Name,
			Year:    form.Year,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		handler.sessions.Put(req.Context(), "flash_success", "Thema wurde erfolgreich bearbeitet.")

		// Redirect to list of events
		http.Redirect(res, req, "/topics/"+topicIDstr+"/events", http.StatusFound)
	}
}
