// The pivot of all HTTP-handlers functions, which is responsible for
// initializing a web handler, consisting of a multiplexer, a database store
// and a session manager. It also contains middleware and singled out HTTP-
// handler functions.

package web

import (
	"context"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/csrf"

	x "github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

const (
	path   = "frontend/templates/"
	layout = "frontend/layout.html"
	css    = "frontend/css/css.html"
)

var (
	// _testing is a flag to skip init function when testing
	_testing = false

	// Parsed HTML-templates to be executed in their respective HTTP-handler
	// functions when needed
	homeTemplate, notFound404Template *template.Template

	// funcMap is a map of custom functions to be used in an HTML-template
	funcMap = template.FuncMap{
		"increment": func(num int) int {
			return num + 1
		},
		"decrement": func(num int) int {
			return num - 1
		},
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

	homeTemplate = template.Must(template.ParseFiles(layout, css, path+"home.html"))
	notFound404Template = template.Must(template.ParseFiles(layout, css, path+"http_not_found.html"))
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

	// Home
	handler.Get("/", handler.Home())

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
	handler.Route("/scores", func(router chi.Router) {
		router.Get("/", scores.List())
		router.Post("/", scores.Filter())
	})

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
	handler.NotFound(handler.NotFound404())

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

		Topics []x.Topic
	}

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute SQL statement to get topics
		topics, err := h.store.GetTopics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err = homeTemplate.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Topics:      topics,
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

// NotFound404 gets called when a non-existing URL has been entered.
func (h *Handler) NotFound404() http.HandlerFunc {

	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if err := notFound404Template.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
