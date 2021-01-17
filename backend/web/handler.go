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

	"github.com/mqrc81/IDPA-Jahreszahlen/backend/jahreszahlen"
)

var (
	// Parsed HTML-templates to be executed in their respective HTTP-handler
	// functions when needed
	Templates = make(map[string]*template.Template)

	// A map of custom functions to be used in an HTML-template
	funcMap = template.FuncMap{
		"is_even": func(num int) bool {
			return num%2 == 0
		},
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
// All HTML-templates get parsed once and added to a map to be executed when
// needed. This is way more efficient than parsing the HTML-templates every
// time a request is sent.
func init() {
	path := "frontend/templates/"
	layout := "frontend/layout.html"
	css := "frontend/css/css.html"

	Templates["home"] = template.Must(template.ParseFiles(layout, css, path+"home.html"))
	Templates["http_not_found"] = template.Must(template.ParseFiles(layout, css, path+"http_not_found.html"))

	Templates["topics_list"] = template.Must(template.ParseFiles(layout, css, path+"topics_list.html"))
	Templates["topics_create"] = template.Must(template.ParseFiles(layout, css, path+"topics_create.html"))
	Templates["topics_edit"] = template.Must(template.ParseFiles(layout, css, path+"topics_edit.html"))
	Templates["topics_show"] = template.Must(template.ParseFiles(layout, css, path+"topics_show.html"))

	Templates["events_list"] = template.Must(template.ParseFiles(layout, css, path+"events_list.html"))
	Templates["events_create"] = template.Must(template.ParseFiles(layout, css, path+"events_create.html"))
	Templates["events_edit"] = template.Must(template.ParseFiles(layout, css, path+"events_edit.html"))

	Templates["scores_list"] = template.Must(template.
		New("layout.html").Funcs(funcMap). // add custom functions to use in HTML-templates
		ParseFiles(layout, css, path+"scores_list.html"))

	Templates["quiz_phase1"] = template.Must(template.ParseFiles(layout, css, path+"quiz_phase1.html"))
	Templates["quiz_phase1_review"] = template.Must(template.ParseFiles(layout, css, path+"quiz_phase1_review.html"))
	Templates["quiz_phase2"] = template.Must(template.ParseFiles(layout, css, path+"quiz_phase2.html"))
	Templates["quiz_phase2_review"] = template.Must(template.ParseFiles(layout, css, path+"quiz_phase2_review.html"))
	Templates["quiz_phase3"] = template.Must(template.ParseFiles(layout, css, path+"quiz_phase3.html"))
	Templates["quiz_phase3_review"] = template.Must(template.ParseFiles(layout, css, path+"quiz_phase3_review.html"))
	Templates["quiz_summary"] = template.Must(template.ParseFiles(layout, css, path+"quiz_summary.html"))

	Templates["users_register"] = template.Must(template.ParseFiles(layout, css, path+"users_register.html"))
	Templates["users_login"] = template.Must(template.ParseFiles(layout, css, path+"users_login.html"))
	Templates["users_profile"] = template.Must(template.ParseFiles(layout, css, path+"users_profile.html"))
	Templates["users_list"] = template.Must(template.ParseFiles(layout, css, path+"users_list.html"))
	Templates["users_edit_password"] = template.Must(template.ParseFiles(layout, css, path+"users_edit_password.html"))
	Templates["users_edit_username"] = template.Must(template.ParseFiles(layout, css, path+"users_edit_username.html"))
	Templates["users_forgot_password"] = template.Must(template.ParseFiles(layout, css,
		path+"users_forgot_password.html"))
	Templates["users_reset_password"] = template.Must(template.ParseFiles(layout, css,
		path+"users_reset_password.html"))
}

// NewHandler initializes HTTP-handlers, including router and middleware.
func NewHandler(store jahreszahlen.Store, sessions *scs.SessionManager, csrfKey []byte) *Handler {
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
		r.Get("/new", topics.Create())
		r.Post("/", topics.CreateStore())
		r.Post("/{topicID}/delete", topics.Delete())
		r.Get("/{topicID}/edit", topics.Edit())
		r.Get("/{topicID}", topics.Show())
	})

	// Events
	handler.Route("/topics/{topicID}/events", func(router chi.Router) {
		router.Get("/", events.List())
		router.Get("/new", events.Create())
		router.Post("/", events.CreateStore())
		router.Post("/{eventID}/delete", events.Delete())
		router.Get("/edit", events.Edit())
		router.Post("/edit", events.EditStore())
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

		router.Post("/verify/email", users.VerifyEmail())
		router.Get("/forgot/password", users.ForgotPassword())
		router.Post("/forgot/password", users.ForgotPasswordSubmit())
		router.Get("/reset/password", users.ResetPassword())
		router.Post("/reset/password", users.ResetPasswordSubmit())

		router.Get("/test/email", users.TestEmail()) // TEMP
	})

	// Handler for when a non-existing URL is called
	handler.NotFound(handler.NotFound404())

	return handler
}

// Handler consists of the chi-multiplexer, a store interface and sessions.
type Handler struct {
	*chi.Mux

	store    jahreszahlen.Store
	sessions *scs.SessionManager
}

// Home is a GET-method that is accessible to anyone.
//
// It displays the home-page.
func (handler *Handler) Home() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topics []jahreszahlen.Topic
	}

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute SQL statement to get topics
		topics, err := handler.store.GetTopics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := Templates["home"].Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topics:      topics,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// withUser is a middleware that replaces the potential user ID with a user object.
func (handler *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Retrieve user ID from session
		var userID int
		userIDinf := handler.sessions.Get(req.Context(), "user_id")
		if userIDinf != nil {
			userID = userIDinf.(int)
		}

		// Execute SQL statement to get user
		user, err := handler.store.GetUser(userID)
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
func (handler *Handler) NotFound404() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}

	return func(res http.ResponseWriter, req *http.Request) {
		if err := Templates["http_not_found"].Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// Abs gets absolute value of an int number (-10 => 10)
func Abs(num int) int {
	if num < 0 {
		return -num
	}
	return num
}
