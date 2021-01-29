// The web handler evolving around events, with HTTP-handler functions
// consisting of "GET"- and "POST"-methods. It utilizes session management and
// database access.

package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/gorilla/csrf"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	eventsListTemplate, eventsCreateTemplate, eventsEditTemplate *template.Template
)

func init() {
	if _testing { // skip initialization of templates when running tests
		return
	}

	eventsListTemplate = template.Must(template.ParseFiles(layout, css, path+"events_list.html"))
	eventsCreateTemplate = template.Must(template.ParseFiles(layout, css, path+"events_create.html"))
	eventsEditTemplate = template.Must(template.ParseFiles(layout, css, path+"events_edit.html"))
}

// EventHandler is the object for handlers to access sessions and database.
type EventHandler struct {
	store    x.Store
	sessions *scs.SessionManager
}

// List is a GET-method that is accessible to any admin. It lists all events,
// sorted by date ascending.
//
// Users can view the events while admins have the ability to edit or delete an
// event, as well as to create a new one.
func (h *EventHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML

		Topic  x.Topic
		Events []x.Event
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		user := req.Context().Value("user")
		if user == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
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
		topic, err := h.store.GetTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err = eventsListTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
			Topic:       topic,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Create is a GET-method that is accessible to any admin.
//
// It displays a form, in which values for a new event can be entered.
func (h *EventHandler) Create() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML

		Topic x.Topic
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(x.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
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
		topic, err := h.store.GetTopic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err = eventsCreateTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
			Topic:       topic,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CreateStore is a POST-method that is accessible to anyone after Create.
//
// It validates the form from Create and redirects to Create in case of an
// invalid input with the corresponding error message. In case of valid form,
// it stores the new event in the database and redirects to the List.
func (h *EventHandler) CreateStore() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		form := EventForm{
			Name:       req.FormValue("name"),
			YearOrDate: req.FormValue("year"),
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Execute SQL statement to create an event
		if err := h.store.CreateEvent(&x.Event{
			TopicID: topicID,
			Name:    form.Name,
			Year:    form.Year,
			Date:    form.Date,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Adds flash message
		h.sessions.Put(req.Context(), "flash_success", "Ereignis wurde erfolgreich erstellt.")

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicIDstr+"/events", http.StatusFound)
	}
}

// Delete is a POST-method that is accessible to any admin after List.
//
// It deletes an event and redirects to List.
func (h *EventHandler) Delete() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve event ID from URL parameters
		topicID := chi.URLParam(req, "topicID")

		// Retrieve event ID from URL parameters
		eventID, _ := strconv.Atoi(chi.URLParam(req, "eventID"))

		// Execute SQL statement to delete an event
		if err := h.store.DeleteEvent(eventID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of topics
		http.Redirect(res, req, "/topics/"+topicID+"/events", http.StatusFound)
	}
}

// EditPrepare is a POST-method that is accessible to any admin.
//
// It creates a form from the event's values, so that values are already filled
// out when editing the event.
func (h *EventHandler) EditPrepare() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from URL parameters
		topicID := chi.URLParam(req, "topicID")
		eventIDstr := chi.URLParam(req, "topicID")
		eventID, _ := strconv.Atoi(eventIDstr)

		// Execute SQL statement to get topic
		event, err := h.store.GetEvent(eventID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create form to add to session, so that values are already filled out
		// when editing a topic
		form := EventForm{
			Name:   event.Name,
			Year:   event.Year,
			Errors: FormErrors{},
		}

		// Add form to session
		h.sessions.Put(req.Context(), "form", form)

		// Redirect to edit-page of topic
		http.Redirect(res, req, "/topics/"+topicID+"/events/"+eventIDstr+"edit", http.StatusFound)
	}
}

// Edit is a GET-method that is accessible to any admin.
//
// It displays a form in which values for modifying the current event can be
// entered.
func (h *EventHandler) Edit() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML

		Event x.Event
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(x.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
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
		event, err := h.store.GetEvent(eventID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err = eventsEditTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
			Event:       event,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// EditStore is a POST-method that is accessible to any admin after Edit.
//
// It validates the form from Edit and redirects to Edit in case of an invalid
// input with the corresponding error message. In case of valid form, it stores
// the topic in the database and redirects to List.
func (h *EventHandler) EditStore() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve topic ID from URL parameters
		topicIDstr := chi.URLParam(req, "topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve values from form
		form := EventForm{
			Name:       req.FormValue("name"),
			YearOrDate: req.FormValue("year"),
		}

		// Validate form
		if !form.Validate() {
			h.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to update event
		if err := h.store.UpdateEvent(&x.Event{
			TopicID: topicID,
			Name:    form.Name,
			Year:    form.Year,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message to session
		h.sessions.Put(req.Context(), "flash_success", "Thema wurde erfolgreich bearbeitet.")

		// Redirect to list of events
		http.Redirect(res, req, "/topics/"+topicIDstr+"/events", http.StatusFound)
	}
}
