package web

// topic_handler.go
// Contains all HTTP-handlers for pages evolving around topics.

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// TopicHandler
// Object for handlers to access sessions and database.
type TopicHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// List
// A GET-method that any admin can call. It lists all topics.
func (handler *TopicHandler) List() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topics     []backend.Topic
		EventCount int
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/topics_list.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute SQL statement to get topics
		topics, err := handler.store.Topics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// TODO Retrieve topic ID from URL

		// TODO Execute SQL statement to get amount of events

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topics:      topics,
			// TODO EventCount: eCount,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Create
// A GET-method that any admin can call. It renders a form, in which values
// for a new topic can be entered.
func (handler *TopicHandler) Create() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/topics_create.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Check if logged in user is admin
		user := req.Context().Value("user")
		if user == nil || !user.(backend.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// redirect back
			handler.sessions.Put(req.Context(), "flash", "NOOOOOPE")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// CreateStore
// A POST-method. It validates the form from Create and redirects to Create in
// case of an invalid input with corresponding error message. In case of valid
// form, it stores the new topic in the database and redirects to List.
func (handler *TopicHandler) CreateStore() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve variables from form (Create)
		startYear, _ := strconv.Atoi(req.FormValue("start_year"))
		endYear, _ := strconv.Atoi(req.FormValue("end_year"))
		form := TopicForm{
			Title:       req.FormValue("title"),
			StartYear:   startYear,
			EndYear:     endYear,
			Description: req.FormValue("description"),
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to create a topic
		if err := handler.store.CreateTopic(&backend.Topic{
			Title:       form.Title,
			StartYear:   form.StartYear,
			EndYear:     form.EndYear,
			Description: form.Description,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Adds flash message
		handler.sessions.Put(req.Context(), "flash", "Thema wurde erfolgreich erstellt.")

		// Redirects to list of topics
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

// Delete
// A POST-method. It deletes a certain topic and redirects to Show.
func (handler *TopicHandler) Delete() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve TopicID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute SQL statement to delete a topic
		if err := handler.store.DeleteTopic(topicID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		// Add flash message
		handler.sessions.Put(req.Context(), "flash", "Thema wurde erfolgreich gelöscht.")

		// Redirect to list of topics
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

// Edit
// A GET-method that any admin can call. It renders a form in which values for
// updating the current topic can be entered.
func (handler *TopicHandler) Edit() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topic  backend.Topic
		Events []backend.Event
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/topics_edit.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute SQL statement to get a topic
		topic, err := handler.store.Topic(topicID)
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

// EditStore
// A POST-method. It validates the form from Edit and redirects to Edit in
// case of an invalid input with corresponding error message. In case of valid
// form, it stores the topic in the database and redirects to Show.
func (handler *TopicHandler) EditStore() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve topic ID from URL
		topicIDstr := req.URL.Query().Get("topicID")
		topicID, _ := strconv.Atoi(topicIDstr)

		// Retrieve values from form (Edit)
		startYear, _ := strconv.Atoi(req.FormValue("start_year"))
		endYear, _ := strconv.Atoi(req.FormValue("end_year"))
		form := TopicForm{
			Title:       req.FormValue("title"),
			StartYear:   startYear,
			EndYear:     endYear,
			Description: req.FormValue("description"),
		}

		// Validate form
		if !form.Validate() {
			handler.sessions.Put(req.Context(), "form", form)
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute SQL statement to update a topic
		if err := handler.store.UpdateTopic(&backend.Topic{
			TopicID:     topicID,
			Title:       form.Title,
			StartYear:   form.StartYear,
			EndYear:     form.EndYear,
			Description: form.Description,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add flash message
		handler.sessions.Put(req.Context(), "flash", "Thema wurde erfolgreich bearbeitet.")

		// Redirect to topic Show
		http.Redirect(res, req, "/topics/"+topicIDstr, http.StatusFound)
	}
}

// Show
// A GET-method that any user can call. It displays details of the topic and
// has the options to play or edit the topics, to edit an event and to create a
// new event.
func (handler *TopicHandler) Show() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topic backend.Topic
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/topics_show.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve TopicID from URL
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute SQL statement to get a topic
		topic, err := handler.store.Topic(topicID)
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
