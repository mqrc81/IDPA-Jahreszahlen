package web

/*
 * Contains all HTTP-handler functions for pages evolving around scores.
 */

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreHandler is the object for handlers to access sessions and database.
type ScoreHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// TODO
// List is a GET-method that any user can call. It lists all scores, ranked by
// points, with the ability to filter scores by topic and/or user.
func (handler *ScoreHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Scores []backend.Score
	}

	// Parse HTML-templates
	tmpl := template.Must(template.New("").ParseFiles(
		"frontend/layout.html",
		"frontend/pages/scores_list.html",
	))

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		userInf := req.Context().Value("user")
		if userInf == nil {
			// If no user is logged in, then redirect back with flash message
			handler.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Benutzer eingeloggt sein, um das Leaderboard zu betrachten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}
		user := userInf.(backend.User)

		var scores []backend.Score

		// Retrieve topic from URL parameters
		topicID := -1
		var err error
		topic := req.URL.Query().Get("topic")
		if len(topic) != 0 {
			topicID, err = strconv.Atoi(topic)
			if err != nil {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
		}

		// Retrieve user from URL parameters
		userFilter := req.URL.Query().Get("user")

		// Retrieve limit from URL parameters
		limit := 15
		show := req.URL.Query().Get("show")
		if len(show) != 0 {
			limit, err = strconv.Atoi(show)
			if err != nil {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
		}

		// Retrieve offset from URL parameters
		offset := 0
		page := req.URL.Query().Get("page")
		if len(page) != 0 {
			offset, err = strconv.Atoi(page)
			if err != nil {
				http.Error(res, err.Error(), http.StatusNotFound)
				return
			}
			offset = (offset - 1) * limit
		}

		if topicID != -1 {
			if strings.ToLower(userFilter) == "me" { // Topic and user specified in URL parameters
				// Execute SQL statement to get scores
				scores_, err := handler.store.GetScoresByTopicAndUser(topicID, user.UserID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				scores = scores_
			} else { // Only topic specified in URL parameters
				// Execute SQL statement to get scores
				scores_, err := handler.store.GetScoresByTopic(topicID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				scores = scores_
			}
		} else if strings.ToLower(userFilter) == "me" { // Only user specified in URL parameters
			// Execute SQL statement to get scores
			scores_, err := handler.store.GetScoresByUser(user.UserID, limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = scores_
		} else { // No topic or user specified in URL parameters
			// Execute SQL statement to get scores
			scores_, err := handler.store.GetScores(limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = scores_
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Scores:      scores,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Store is a POST-method. It stores the new score in the database and
// redirects to List.
func (handler *ScoreHandler) Store() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from URL
		topicID, _ := strconv.Atoi(req.URL.Query().Get("topic_id"))
		userID, _ := strconv.Atoi(req.URL.Query().Get("user_id"))
		points, _ := strconv.Atoi(req.URL.Query().Get("points"))
		date := time.Now().Format("2006-01-02")

		// Execute SQL statement to create a score
		if err := handler.store.CreateScore(&backend.Score{
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
