// The web handler evolving around scores, with HTTP-handler functions
// consisting of "GET"- and "POST"-methods. It utilizes session management and
// database access.

package web

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/alexedwards/scs/v2"
	"github.com/gorilla/csrf"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// Parsed HTML-templates to be executed in their respective HTTP-handler
	// functions when needed
	scoresListTemplate *template.Template
)

const (
	showDefault = 10
)

// init gets initialized with the package.
//
// All HTML-templates get parsed once to be executed when needed. This is way
// more efficient than parsing the HTML-templates with every request.
func init() {
	if _testing { // skip initialization of templates when running tests
		return
	}

	scoresListTemplate = template.Must(template.ParseFiles(layout, templatePath+"scores_list.html"))
}

// ScoreHandler is the object for handlers to access sessions and database.
type ScoreHandler struct {
	store    x.Store
	sessions *scs.SessionManager
}

// List is a GET-method that is accessible to any user.
//
// It lists all scores and displays it as a leaderboard table, ranked by
// points, with the ability to filter whilst typing in the search bar, as well
// as to choose how many entries are shown at a time and navigate to the
// previous or next page.
//
// The leaderboard contains of a rank, name of user, name of topic, date and
// points of a score.
func (h *ScoreHandler) List() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
		CSRF template.HTML

		Leaderboard []leaderboardRow
		Show        int  // amount of scores to be shown per page
		Page        int  // current page
		PrevPage    bool // false if current page is 1
		NextPage    bool // false if scores array ends on current page
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Check if a user is logged in
		userInf := req.Context().Value("user")
		if userInf == nil {
			// If no user is logged in, then redirect back with flash message
			h.sessions.Put(req.Context(), "flash_error", "Unzureichende Berechtigung. "+
				"Sie müssen als Benutzer eingeloggt sein, um das Leaderboard zu betrachten.")
			http.Redirect(res, req, req.Referer(), http.StatusFound)
			return
		}

		// Retrieve values for filtering the leaderboard from sessions
		show := h.sessions.GetInt(req.Context(), "show")
		if show == 0 { // 'show' is 0 in case of no 'show' in the session
			show = showDefault
		}
		page := h.sessions.GetInt(req.Context(), "page")
		if page == 0 { // 'page' is 0 in case of no 'page' in the session
			page = 1
		}

		// Execute SQL statement to get scores
		scores, err := h.store.GetScores()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Create table of scores for the current page & show size
		// Example: page = 3, show = 15 => scores[30:45] (ranks 31-45)
		// However, if len(scores) = 33 => scores[30:33] (ranks 31-33)
		leaderboard := createLeaderboardRows(scores, show, page)

		// Execute HTML-templates with data
		if err = scoresListTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			CSRF:        csrf.TemplateField(req),
			Leaderboard: leaderboard,
			Show:        show,
			Page:        page,
			PrevPage:    page-1 > 0,
			NextPage:    page-1 >= (len(scores)-1)/show,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Filter filters show
func (h *ScoreHandler) Filter() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve values from form
		show, _ := strconv.Atoi(req.FormValue("show_filter"))
		pageFilter := req.FormValue("page_filter")
		page, _ := strconv.Atoi(req.FormValue("current_page"))

		switch pageFilter {
		case "prev":
			page--
		case "next":
			page++
		}

		// Add filtering options to sessions
		h.sessions.Put(req.Context(), "show", show)
		h.sessions.Put(req.Context(), "page", page)

		// Redirect to leaderboard
		http.Redirect(res, req, "/scores", http.StatusFound)
	}
}

// leaderboardRow represents a row of the leaderboard.
type leaderboardRow struct {
	Rank      int
	UserName  string
	TopicName string
	Date      string
	Points    int
}

// createLeaderboardRows generates all rows of the leaderboard. 'show'
// indicates the amount of scores (scores[n:n+show]) and 'page' indicates
// the offset of the range (start=show*(page-1) -> scores[start:start+show]).
func createLeaderboardRows(scores []x.Score, show int, page int) []leaderboardRow {
	var leaderboard []leaderboardRow

	start := show * (page - 1)
	end := show * page
	for i := start; i < len(scores) && i < end; i++ {
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
