package web

// score_handler.go
// Contains all HTTP-handlers for pages evolving around scores.

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreHandler
// Object for handlers to access sessions and database.
type ScoreHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// List
// A GET-method that any user can call. It lists all scores, ranked by points,
// with the ability to filter scores by topic and/or user.
func (handler *ScoreHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Scores []backend.Score
	}

	// Parse HTML-templates
	tmpl := template.Must(template.New("").Funcs(FuncMap).ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/scores_list.html"))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		userInf := req.Context().Value("user")
		if userInf == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. " +
				"Sie müssen als Benutzer eingeloggt sein, um das Leaderboard zu betrachten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}
		user := userInf.(backend.User)

		var ss []backend.Score

		// Retrieve topic from URL parameters
		topicID := -1
		topic := req.URL.Query().Get("topic")
		if len(topic) != 0 {
			topicID, _ = strconv.Atoi(topic)
		}

		// Retrieve user from URL parameters
		userFilter := req.URL.Query().Get("user")

		// Retrieve limit from URL parameters
		limit := 15
		show := req.URL.Query().Get("show")
		if len(show) != 0 {
			limit, _ = strconv.Atoi(show)
		}

		// Retrieve offset from URL parameters
		offset := 0
		page := req.URL.Query().Get("page")
		if len(page) != 0 {
			offset, _ = strconv.Atoi(page)
			offset = (offset - 1) * limit
		}

		if topicID != -1 {
			if strings.ToLower(userFilter) == "me" { // Topic and user specified in URL parameters
				// Execute SQL statement to get scores
				scores, err := handler.store.ScoresByTopicAndUser(topicID, user.UserID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				ss = scores
			} else { // Only topic specified in URL parameters
				// Execute SQL statement to get scores
				scores, err := handler.store.ScoresByTopic(topicID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				ss = scores
			}
		} else if strings.ToLower(userFilter) == "me" { // Only user specified in URL parameters
			// Execute SQL statement to get scores
			scores, err := handler.store.ScoresByUser(user.UserID, limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			ss = scores
		} else { // No topic or user specified in URL parameters
			// Execute SQL statement to get scores
			scores, err := handler.store.Scores(limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			ss = scores
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Scores:      ss,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Store
// A POST-method. It stores the new score in the database and redirects to List.
func (handler *ScoreHandler) Store() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		topicID, _ := strconv.Atoi(req.URL.Query().Get("topic_id"))
		userID, _ := strconv.Atoi(req.URL.Query().Get("user_id"))
		points, _ := strconv.Atoi(req.URL.Query().Get("points"))
		date := time.Now().Format("2006-01-02")

		// Execute SQL statement to create a score
		if err := handler.store.CreateScore(&backend.Score{
			ScoreID: 0,
			TopicID: topicID,
			UserID:  userID,
			Points:  points,
			Date:    date,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
