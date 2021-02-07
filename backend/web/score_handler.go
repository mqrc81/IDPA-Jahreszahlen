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

	scoresListTemplate = template.Must(template.New("layout.html").
		Funcs(template.FuncMap{ // Add custom HTML-template-function to increment a number
			"increment": func(num int) int {
				return num + 1
			},
			"decrement": func(num int) int {
				return num - 1
			},
		}).
		ParseFiles(layout, templatePath+"scores_list.html"))
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

		Show     int  // amount of scores shown
		ShowFrom int  // first score's rank
		ShowTo   int  // last score's rank
		ShowOf   int  // total amount of scores
		ShowAll  bool // whether all scores are shown

		Page         int   // current page
		Pages        []int // range of pages to be able to navigate to
		PagePrevious bool  // whether there's a previous page
		PageNext     bool  // whether there's a next page
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

		// Execute SQL statement to get scores
		scores, err := h.store.GetScores()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Retrieve values from URL query for filtering the leaderboard by
		// indicating the amount of scores to be shown and with which offset
		showFilter := req.URL.Query().Get("show")
		pageFilter := req.URL.Query().Get("page")
		show, page := inspectFilters(showFilter, pageFilter, len(scores))

		// Create table of scores for the current page & show size
		// Example: page = 3, show = 15 => scores[30:45] (ranks 31-45)
		// However, if len(scores) = 33 => scores[30:33] (ranks 31-33)
		leaderboard := createLeaderboardRows(scores, show, page)

		// Page numbers to be shown below the leaderboard in order to navigate
		// to different pages
		pages := createPages(show, page, len(scores))

		// Execute HTML-templates with data
		if err = scoresListTemplate.Execute(res, data{
			SessionData:  GetSessionData(h.sessions, req.Context()),
			CSRF:         csrf.TemplateField(req),
			Leaderboard:  leaderboard,
			Show:         show,
			ShowFrom:     leaderboard[0].Rank,
			ShowTo:       leaderboard[len(leaderboard)-1].Rank,
			ShowOf:       len(scores),
			ShowAll:      show == len(scores),
			Page:         page,
			Pages:        pages,
			PagePrevious: page != pages[0],
			PageNext:     page != pages[len(pages)-1],
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
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

// inspectFilters examines the filters from the URL query, checks for all
// possible cases and returns the amount of scores to be shown and the
// new page.
// (Tested in test_handler.go)
func inspectFilters(showFilter string, pageFilter string, scoresCount int) (int, int) {

	// Check for invalid URL query
	show, _ := strconv.Atoi(showFilter)
	page, _ := strconv.Atoi(pageFilter)

	// Inspect 'show'
	// 'show' = 0 if no 'show' URL query was found
	// 'show' may only be 10, 25, 50 or all
	if show == 0 {
		show = showDefault
	} else if show < 0 {
		show = scoresCount
	} else if show != 10 && show != 25 && show != 50 {
		if show < 20 {
			show = 10
		} else if show < 40 {
			show = 25
		} else {
			show = 50
		}
	}

	// Inspect 'page'
	// Catch leaderboard rows out of bounds (user is on page 2, showing 10
	// rows, total scores are 22 and user chooses to show 25 rows => page drops
	// from 2 to 1)
	if page < 1 || show == scoresCount {
		page = 1
	} else if scoresCount < (page-1)*show {
		page = ((scoresCount - 1) / show) + 1
	}

	return show, page
}

// createPages generates the pages to which the user can navigate to.
// Example 1: 23 scores, showing 11-20 => [< 1 '2' 3 >]
// Example 2: 20 scores, showing 11-20 => [< 1 '2']
// Example 3: 50 scores, showing 41-50 => [< 3 4 '5']
// (Tested in test_handler.go)
func createPages(show int, page int, scoresCount int) []int {
	var pages []int

	if page == 1 { // case ['1' 2 3 >]
		for i := 0; i <= 2; i++ {
			if scoresCount > (page+i-1)*show {
				pages = append(pages, page+i)
			}
		}
	} else if scoresCount <= show*page { // case [< 2 3 '4']
		for i := 2; i >= 0; i-- {
			if page > i {
				pages = append(pages, page-i)
			}
		}
	} else {
		for i := -1; i <= 1; i++ { // case [< 1 '2' 3 >]
			if scoresCount > (page+i-1)*show {
				pages = append(pages, page+i)
			}
		}
	}

	return pages
}
