package web

/*
 * Contains all HTTP-handler functions for pages evolving around scores.
 */

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

// ScoreHandler is the object for handlers to access sessions and database.
type ScoreHandler struct {
	store    backend.Store
	sessions *scs.SessionManager
}

// List is a GET-method that any user can call. It lists all scores, ranked by
// points, with the ability to filter scores by topic and/or user.
func (handler *ScoreHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Leaderboard []leaderboardRow

		Topic int
		User  string
		Page  int
		Show  int
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

		// Retrieve URL query parameters for filtering the leaderboard from URL
		topicID, _ := strconv.Atoi(req.URL.Query().Get("topic")) // if topic query is empty or not a number, topicID = 0
		userFilter := strings.ToLower(req.URL.Query().Get("user"))
		limit, _ := strconv.Atoi(req.URL.Query().Get("show")) // if show query is empty or not a number, limit = 25
		if limit == 0 {
			limit = 25
		}
		page, _ := strconv.Atoi(req.URL.Query().Get("page")) // if page is empty or not a number, offset = 0
		offset := (page - 1) * limit                         // if "/scores?page=3&show=15", start at score 31

		if topicID != 0 {
			if userFilter == "me" { // Topic and user specified in URL parameters
				// Execute SQL statement to get scores
				scoresByTopicAndUser, err := handler.store.GetScoresByTopicAndUser(topicID, user.UserID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				scores = scoresByTopicAndUser
			} else { // Only topic specified in URL parameters
				// Execute SQL statement to get scores
				scoresByTopic, err := handler.store.GetScoresByTopic(topicID, limit, offset)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				scores = scoresByTopic
			}
		} else if userFilter == "me" { // Only user specified in URL parameters
			// Execute SQL statement to get scores
			scoresByUser, err := handler.store.GetScoresByUser(user.UserID, limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = scoresByUser
		} else { // No topic or user specified in URL parameters
			// Execute SQL statement to get scores
			scoresAll, err := handler.store.GetScores(limit, offset)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = scoresAll
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Leaderboard: createLeaderboardRows(scores, offset),
			Topic:       topicID,
			User:        userFilter,
			Page:        page,
			Show:        limit,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// leaderboardRow represents a row of the leaderboard
type leaderboardRow struct {
	Rank      int
	UserName  string
	TopicName string
	Date      string
	Points    int
}

// createLeaderboardRows generates all rows of the leaderboard
func createLeaderboardRows(scores []backend.Score, offset int) []leaderboardRow {
	var leaderboard []leaderboardRow

	for index, score := range scores {
		leaderboard = append(leaderboard, leaderboardRow{
			Rank:      index + offset + 1,
			UserName:  score.UserName,
			TopicName: score.TopicName,
			Date:      score.Date,
			Points:    score.Points,
		})
	}

	return leaderboard
}
