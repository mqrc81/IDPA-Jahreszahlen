// The pivot of all HTTP-handlers functions, which is responsible for
// initializing a web handler, consisting of a multiplexer, a database store
// and a session manager, as well as serving static files. It also contains
// middleware and some general HTTP-handler functions.

package web

import (
	"context"
	"html/template"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

const (
	staticPath   = "frontend/static"
	templatePath = "frontend/html/templates/"
	layout       = "frontend/html/layout.html"

	topicsURL   = "/topics"
	scoresURL   = "/scores"
	profileURL  = "/users/profile"
	loginURL    = "/users/login"
	registerURL = "/users/login"
	usersURL    = "/users"
)

var (
	// _testing is a flag to skip init function when testing
	_testing = false

	// Parsed HTML-templates to be executed in their respective HTTP-handler
	// functions when needed
	homeTemplate, http404Template, http405Template *template.Template

	// Possible search result matches from the navigation bar search input
	searchKeywords = map[string]string{
		"quiz":       topicsURL,
		"thema":      topicsURL,
		"themen":     topicsURL,
		"ereignisse": topicsURL,
		"spielen":    topicsURL,
		"ereignis":   topicsURL,

		"leaderboard": scoresURL,
		"ranking":     scoresURL,
		"tabelle":     scoresURL,
		"resultat":    scoresURL,
		"resultate":   scoresURL,
		"punkte":      scoresURL,

		"account":      profileURL,
		"konto":        profileURL,
		"profil":       profileURL,
		"passwort":     profileURL,
		"mail":         profileURL,
		"email":        profileURL,
		"e-mail":       profileURL,
		"username":     profileURL,
		"benutzername": profileURL,

		"login":     loginURL,
		"einloggen": loginURL,
		"anmelden":  loginURL,

		"register":     registerURL,
		"registrieren": registerURL,

		"user":      usersURL,
		"users":     usersURL,
		"benutzer":  usersURL,
		"verwalten": usersURL,
		"befördern": usersURL,
		"admin":     usersURL,
	}
)

// init gets initialized with the package.
//
// All HTML-templates get parsed once to be executed when needed. This is way
// more efficient than parsing the HTML-templates with every request.
func init() {
	if _testing { // skip initialization of templates when running tests
		return
	}

	homeTemplate = template.Must(template.New("layout.html").
		Funcs(template.FuncMap{ // Add custom HTML-template-function to increment a number
			"increment": func(num int) int {
				return num + 1
			},
		}).
		ParseFiles(layout, templatePath+"home.html"))
	http404Template = template.Must(template.ParseFiles(layout, templatePath+"http_not_found.html"))
	http405Template = template.Must(template.ParseFiles(layout, templatePath+"http_method_not_allowed.html"))
}

// NewHandler initializes HTTP-handlers, including router and middleware.
func NewHandler(store x.Store, sessions *scs.SessionManager, csrfKey []byte) *Handler {
	handler := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	topics := TopicHandler{store: store, sessions: sessions}
	events := EventHandler{store: store, sessions: sessions}
	scores := ScoreHandler{store: store, sessions: sessions}
	quiz := QuizHandler{store: store, sessions: sessions}
	users := UserHandler{store: store, sessions: sessions}

	// Use middleware
	handler.Use(middleware.Logger)
	handler.Use(csrf.Protect(csrfKey, csrf.Secure(false)))
	handler.Use(sessions.LoadAndSave)
	handler.Use(handler.withUser)

	// Serve static files
	handler.fileServer("/"+staticPath+"/", http.Dir(staticPath))

	// Home
	handler.Get("/", handler.Home())
	handler.Get("/search", handler.Search())

	// Topics
	handler.Route("/topics", func(r chi.Router) {
		r.Get("/", topics.List())
		r.Get("/{topicID}", topics.Show())
		r.Get("/new", topics.Create())
		r.Post("/", topics.CreateStore())
		r.Post("/{topicID}/delete", topics.Delete())
		r.Get("/{topicID}/edit", topics.Edit())
		r.Post("/{topicID}/edit", topics.EditStore())
	})

	// Events
	handler.Route("/topics/{topicID}/events", func(router chi.Router) {
		router.Get("/", events.List())
		router.Get("/new", events.Create())
		router.Post("/", events.CreateStore())
		router.Post("/{eventID}/delete", events.Delete())
		router.Get("/{eventID}/edit", events.Edit())
		router.Post("/{eventID}/edit", events.EditStore())
	})

	// Quiz
	handler.Route("/topics/{topicID}/quiz", func(router chi.Router) {
		router.Get("/1", quiz.Phase1())
		router.Post("/1", quiz.Phase1Submit())
		router.Get("/1/review", quiz.Phase1Review())
		router.Post("/1/review", quiz.Phase2Prepare())
		router.Get("/2", quiz.Phase2())
		router.Post("/2", quiz.Phase2Submit())
		router.Get("/2/review", quiz.Phase2Review())
		router.Post("/2/review", quiz.Phase3Prepare())
		router.Get("/3", quiz.Phase3())
		router.Post("/3", quiz.Phase3Submit())
		router.Get("/3/review", quiz.Phase3Review())
		router.Get("/summary", quiz.Summary())
	})

	// Scores
	handler.Get("/scores", scores.List())

	// Users
	handler.Route("/users", func(router chi.Router) {
		router.Get("/register", users.Register())
		router.Post("/register", users.RegisterSubmit())
		router.Get("/login", users.Login())
		router.Post("/login", users.LoginSubmit())
		router.Get("/logout", users.Logout())
		router.Get("/profile", users.Profile())
		router.Get("/", users.List())
		router.Post("/{userID}/delete", users.Delete())
		router.Post("/{userID}/promote", users.Promote())

		router.Get("/edit/username", users.EditUsername())
		router.Post("/edit/username", users.EditUsernameSubmit())
		router.Get("/edit/email", users.EditEmail())
		router.Post("/edit/email", users.EditEmailSubmit())
		router.Get("/edit/password", users.EditPassword())
		router.Post("/edit/password", users.EditPasswordSubmit())

		router.Get("/verify/email", users.VerifyEmail())
		router.Post("/resend/email", users.ResendVerifyEmail())
		router.Get("/forgot/password", users.ForgotPassword())
		router.Post("/forgot/password", users.ForgotPasswordSubmit())
		router.Get("/reset/password", users.ResetPassword())
		router.Post("/reset/password", users.ResetPasswordSubmit())
	})

	// Handler for when a non-existing URL is called
	handler.NotFound(handler.HTTP404())
	handler.MethodNotAllowed(handler.HTTP405())

	return handler
}

// Handler consists of the chi-multiplexer, a store interface and sessions.
type Handler struct {
	*chi.Mux

	store    x.Store
	sessions *scs.SessionManager
}

// Home is a GET-method that is accessible to anyone.
//
// It displays the home-page.
func (h *Handler) Home() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topics             []x.Topic
		UsersCount         int
		EventsCount        int
		ScoresCount        int
		ScoresCountMonthly int
	}

	return func(res http.ResponseWriter, req *http.Request) {

		// Execute SQL statement to get topics
		topics, err := h.store.GetTopics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Sort topics by amount of scores in descending order
		sort.Slice(topics, func(n1, n2 int) bool {
			return topics[n1].ScoresCount > topics[n2].ScoresCount
		})
		topics = topics[:min(len(topics), 5)] // only use the 5 topics with the highest amount of scores

		// Execute SQL statement to get amount of users
		usersCount, err := h.store.CountUsers()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get amount of events
		eventsCount, err := h.store.CountEvents()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get amount of scores
		scoresCount, err := h.store.CountScores()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get amount of scores
		scoresCountMonthly, err := h.store.CountScoresByDate(time.Now().AddDate(0, -1, 0), time.Now())
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err = homeTemplate.Execute(res, data{
			SessionData:        GetSessionData(h.sessions, req.Context()),
			Topics:             topics,
			UsersCount:         usersCount,
			EventsCount:        eventsCount,
			ScoresCount:        scoresCount,
			ScoresCountMonthly: scoresCountMonthly,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Search is a GET-method that is accessible to anyone.
//
// It examines the search-query in the navigation bar and redirects user to a
// fitting handler, if any.
func (h *Handler) Search() http.HandlerFunc {

	return func(res http.ResponseWriter, req *http.Request) {

		// Retrieve search result from form
		searchQuery := req.URL.Query().Get("search")
		searchQueries := strings.Split(strings.ToLower(searchQuery), " ")

		// Loop through possible search results to get redirected
		topics, err := h.store.GetTopics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Add topics to search results
		searchResults := searchKeywords
		for _, topic := range topics {
			topicSplit := strings.Split(strings.ToLower(topic.Name), " ")
			for _, t := range topicSplit {
				if t != "der" && t != "die" && t != "das" {
					searchResults[t] = "/topics/" + strconv.Itoa(topic.TopicID)
				}
			}
		}

		// Loop through possible search results
		for _, search := range searchQueries {
			if searchResults[search] != "" { // redirect in case of match
				http.Redirect(res, req, searchResults[search], http.StatusSeeOther)
				return
			}
		}

		// Search query didn't find a match
		h.sessions.Put(req.Context(), "flash_info",
			"Es wurde kein Suchergebnis gefunden. Versuchen Sie es genauer und in ganzen Worten.")
		http.Redirect(res, req, req.Referer(), http.StatusSeeOther)
	}
}

// HTTP404 gets called when a non-existing URL has been entered.
func (h *Handler) HTTP404() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if err := http404Template.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// HTTP405 gets called when a forbidden method gets called to an existing URL.
//
// Example: /users/1/delete as GET method, even though it should be POST.
func (h *Handler) HTTP405() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if err := http405Template.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// withUser is a middleware that replaces the potential user ID with a user object.
func (h *Handler) withUser(next http.Handler) http.Handler {

	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Retrieve user ID from session
		var userID int
		userIDinf := h.sessions.Get(req.Context(), "user_id")
		if userIDinf != nil {
			userID = userIDinf.(int)
		}

		// Execute SQL statement to get user
		user, err := h.store.GetUser(userID)
		if err != nil {
			// No user in session => continue to HTTP-handler
			next.ServeHTTP(res, req)
			return
		}

		// Add the user logged in to the session
		ctx := context.WithValue(req.Context(), "user", user)

		// Serve HTTP with response-writer and request
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}

// fileServer conveniently sets up a http.FileServer handler to serve static
// files, such as CSS, images and JavaScript.
func (h *Handler) fileServer(path string, dir http.FileSystem) {

	// URL mustn't contain variables '{}' or wildcards '*'
	if strings.ContainsAny(path, "{}*") { // URL parameters can be defined as such ('/foo/{bar}/*/foobar')
		log.Fatal("URL parameters not permitted")
	}

	// Modify URL to not end on '/'
	if path != "/" && path[len(path)-1] != '/' {
		h.Get(path, http.RedirectHandler(path+"/", http.StatusMovedPermanently).ServeHTTP)
		path += "/"
	}
	path += "*"

	// HTTP-handler that serves static files with every request
	h.Get(path, func(res http.ResponseWriter, req *http.Request) {
		ctx := chi.RouteContext(req.Context())
		pathPrefix := strings.TrimSuffix(ctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(dir))
		fs.ServeHTTP(res, req)
	})
}

// min returns the smallest out of all the numbers.
func min(nums ...int) int {

	if len(nums) == 0 {
		return 0
	}

	minNumber := nums[0]
	for _, num := range nums {
		if num < minNumber {
			minNumber = num
		}
	}

	return minNumber
}
