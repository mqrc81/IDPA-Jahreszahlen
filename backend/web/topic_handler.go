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
// A GET-method. It lists all topics.
func (handler *TopicHandler) List() http.HandlerFunc {
	// Data to pass to HTML-pages
	type data struct {
		SessionData

		Topics     []backend.Topic
	}

	// Parse HTML-pages
	tmpl := template.Must(template.ParseFiles(
		"frontend/pages/layout.html",
		"frontend/pages/topics_list.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute SQL statement to get topics
		topics, err := handler.store.Topics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-pages with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topics:      topics,
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
	// Data to pass to HTML-pages
	type data struct {
		SessionData
	}

	// Parse HTML-pages
	tmpl := template.Must(template.ParseFiles(
		"frontend/pages/layout.html",
		"frontend/pages/topics_create.html"))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(backend.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. " +
				"Sie müssen als Admin eingeloggt sein, um ein neues Thema zu erstellen.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Execute HTML-pages with data
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
		handler.sessions.Put(req.Context(), "flash_success", "Thema wurde erfolgreich erstellt.")

		// Redirects to list of topics
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

// Delete
// A POST-method. It deletes a certain topic and redirects to Show.
func (handler *TopicHandler) Delete() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve TopicID from URL parameters
		topicID, _ := strconv.Atoi(chi.URLParam(req, "topicID"))

		// Execute SQL statement to delete a topic
		if err := handler.store.DeleteTopic(topicID); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		// Add flash message
		handler.sessions.Put(req.Context(), "flash_success", "Thema wurde erfolgreich gelöscht.")

		// Redirect to list of topics
		http.Redirect(res, req, "/topics", http.StatusFound)
	}
}

// Edit
// A GET-method that any admin can call. It renders a form in which values for
// updating the current topic can be entered.
func (handler *TopicHandler) Edit() http.HandlerFunc {
	// Data to pass to HTML-pages
	type data struct {
		SessionData

		Topic  backend.Topic
		Events []backend.Event
	}

	// Parse HTML-pages
	tmpl := template.Must(template.ParseFiles(
		"frontend/pages/layout.html",
		"frontend/pages/topics_edit.html"))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if an admin is logged in
		user := req.Context().Value("user")
		if user == nil || !user.(backend.User).Admin {
			// If no user is logged in or logged in user isn't an admin,
			// then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error",
				"Unzureichende Berechtigung. Sie müssen als Admin eingeloggt sein, um ein Thema zu bearbeiten.")
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
		topic, err := handler.store.Topic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-pages with data
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
		handler.sessions.Put(req.Context(), "flash_success", "Thema wurde erfolgreich bearbeitet.")

		// Redirect to topic Show
		http.Redirect(res, req, "/topics/"+topicIDstr, http.StatusFound)
	}
}

// Show
// A GET-method. It displays details of the topic with the options to play the
// quiz, to edit the topic and to edit the events.
func (handler *TopicHandler) Show() http.HandlerFunc {
	// Data to pass to HTML-pages
	type data struct {
		SessionData

		Topic backend.Topic
	}

	// Parse HTML-pages
	tmpl := template.Must(template.ParseFiles(
		"frontend/pages/layout.html",
		"frontend/pages/topics_show.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve TopicID from URL parameters
		topicID, err := strconv.Atoi(chi.URLParam(req, "topicID"))
		if err != nil {
			http.Error(res, err.Error(), http.StatusNotFound)
			return
		}

		// Execute SQL statement to get a topic
		topic, err := handler.store.Topic(topicID)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-pages with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topic:       topic,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
