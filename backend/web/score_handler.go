// The web handler evolving around scores, with HTTP-handler functions
// consisting of "GET"- and "POST"-methods. It utilizes session management and
// database access.

package web

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/alexedwards/scs/v2"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/jahreszahlen"
)

// ScoreHandler is the object for handlers to access sessions and database.
type ScoreHandler struct {
	store    jahreszahlen.Store
	sessions *scs.SessionManager
}

// List is a GET-method that is accessible to any user.
//
// It lists all scores and displays it as a leaderboard table, ranked by
// points, with the ability to limit scores to a single topic and to only the
// active user, as well as to choose how many entries are shown at a time and
// switch move to the previous or next page.
//
// The leaderboard contains of a rank, name of user, name of topic, date and
// points of a score.
func (handler *ScoreHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Leaderboard []leaderboardRow

		Topic        int
		User         bool
		Page         int
		Show         int
		NextPage     bool
		PreviousPage bool
	}

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
		user := userInf.(jahreszahlen.User)

		var scores []jahreszahlen.Score

		// Retrieve URL query parameters for filtering the leaderboard from URL
		topicID, _ := strconv.Atoi(req.URL.Query().Get("topic"))   // if topic query is empty or invalid -> topicID = 0
		userFilter := strings.ToLower(req.URL.Query().Get("user")) // usually "me" or empty
		show, _ := strconv.Atoi(req.URL.Query().Get("show"))       // if show query is empty or invalid -> show = 25
		if show == 0 {
			show = 25
		}
		page, err := strconv.Atoi(req.URL.Query().Get("page"))
		if err != nil || page < 1 {
			page = 1
		}
		allUsers := userFilter != "me"

		if topicID != 0 {
			if !allUsers { // Topic and user specified in URL parameters
				// Execute SQL statement to get scores
				scoresByTopicAndUser, err := handler.store.GetScoresByTopicAndUser(topicID, user.UserID)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				scores = scoresByTopicAndUser
			} else { // Only topic specified in URL parameters
				// Execute SQL statement to get scores
				scoresByTopic, err := handler.store.GetScoresByTopic(topicID)
				if err != nil {
					http.Error(res, err.Error(), http.StatusInternalServerError)
					return
				}
				scores = scoresByTopic
			}
		} else if !allUsers { // Only user specified in URL parameters
			// Execute SQL statement to get scores
			scoresByUser, err := handler.store.GetScoresByUser(user.UserID)
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = scoresByUser
		} else { // Neither topic nor user specified in URL parameters
			// Execute SQL statement to get scores
			scoresAll, err := handler.store.GetScores()
			if err != nil {
				http.Error(res, err.Error(), http.StatusInternalServerError)
				return
			}
			scores = scoresAll
		}

		// Adjust page, if the leaderboard rows would be out of bounds of
		// scores array
		// Example: show=10, page=5, len(scores)=13 => page=2 => scores[10:20]
		// => scores[10:13] (ranks 11-13) gets shown (in bounds of array)
		page -= (page*show - (len(scores) + 1)) / show

		// Create table of scores for the current page & show size
		// Example: page=3, show=15 => scores[30:45] (ranks 31-45)
		// However, if len(scores)=33 => scores[30:33] (ranks 31-33)
		leaderboard := createLeaderboardRows(scores, show, page)

		// Execute HTML-templates with data
		if err := Templates["scores_list"].Execute(res, data{
			SessionData:  GetSessionData(handler.sessions, req.Context()),
			Leaderboard:  leaderboard,
			Topic:        topicID,
			User:         allUsers,
			Page:         page,
			Show:         show,
			PreviousPage: page != 1,
			NextPage:     page-1 < (len(scores)-1)/show,
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
func createLeaderboardRows(scores []jahreszahlen.Score, show int, page int) []leaderboardRow {
	var leaderboard []leaderboardRow
	for i := show * (page - 1); i < len(scores) && i < show*page; i++ {
		leaderboard = append(leaderboard, leaderboardRow{
			Rank:      i + 1,
			UserName:  scores[i].UserName,
			TopicName: scores[i].TopicName,
			Date:      scores[i].Date.Format("02.01.06"), // date formatted as 'dd.mm.yy'
			Points:    scores[i].Points,
		})
	}

	return leaderboard
}
