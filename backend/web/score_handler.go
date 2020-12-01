package web

/*
 * score_handler.go contains HTTP-handler functions for scores
 */

import (
	"html/template"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * ScoreHandler handles sessions, CSRF-protection and database access for scores
 */
type ScoreHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

/*
 * List is a GET method that lists scores of a topic sorted by points
 */
func (h *ScoreHandler) List() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		Scores []backend.Score
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Funcs(FuncMap).Parse(scoresListHTML))
	return func(res http.ResponseWriter, req *http.Request) {
		var ss []backend.Score

		// Retrieve topic from URL parameters
		topicID := -1
		topic := req.URL.Query().Get("topic")
		if len(topic) != 0 {
			topicID, _ = strconv.Atoi(topic)
		}

		// Retrieve topic from URL parameters
		userID := -1
		user := req.URL.Query().Get("topic")
		if len(user) != 0 {
			userID, _ = strconv.Atoi(topic)
		}

		// Retrieve limit from URL parameters
		limit := 25
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
			if userID != -1 { // Topic and User specified in URL parameters
				// Execute SQL statement and return scores
				scores, err := h.store.ScoresByTopicAndUser(topicID, userID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				ss = scores
			} else { // Topic specified in URL parameters
				// Execute SQL statement and return scores
				scores, err := h.store.ScoresByTopic(topicID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				ss = scores
			}
		} else if userID != -1 { // User specified in URL parameters
			// Execute SQL statement and return scores
			scores, err := h.store.ScoresByUser(userID, limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			ss = scores
		} else { // No Topic or User specified in URL parameters
			// Execute SQL statement and return scores
			scores, err := h.store.Scores(limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			ss = scores
		}

		// Execute HTML-template
		if err := tmpl.Execute(res, data{Scores: ss}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Store is a POST method that stores new score created
 */
func (h *ScoreHandler) Store() http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		// Retrieve values from form
		topicID, _ := strconv.Atoi(req.URL.Query().Get("topic_id"))
		userID, _ := strconv.Atoi(req.URL.Query().Get("user_id"))
		points, _ := strconv.Atoi(req.URL.Query().Get("points"))
		date := time.Now().Format("2006-01-02")

		// Execute SQL statement
		if err := h.store.CreateScore(&backend.Score{
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
