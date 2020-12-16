package web

// handler.go
// Contains HTTP-router and all HTTP-handlers

import (
	"context"
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	// FuncMap
	// A map that stores functions to use in HTML-template
	FuncMap = template.FuncMap{
		// ranks scores
		"rank": func(num int, page int, limit int) int {
			return (page-1)*limit + num + 1
		},
		// increments number by 1
		"increment": func(num int) int {
			return num + 1
		},
	}
)

// NewHandler
// Initializes HTTP-handlers, including router and middleware
func NewHandler(store backend.Store, sessions *scs.SessionManager) *Handler {
	handler := &Handler{
		Mux:      chi.NewMux(),
		store:    store,
		sessions: sessions,
	}

	topics := TopicHandler{store: store, sessions: sessions}
	events := EventHandler{store: store, sessions: sessions}
	scores := ScoreHandler{store: store, sessions: sessions}
	play := PlayHandler{store: store, sessions: sessions}
	users := UserHandler{store: store, sessions: sessions}

	// Use middleware
	handler.Use(middleware.Logger)
	handler.Use(sessions.LoadAndSave)
	handler.Use(handler.withUser)

	// Home
	handler.Get("/", handler.Home())
	handler.Get("/about", handler.About())

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
		router.Post("/", events.Store())
		router.Post("/{eventID}/delete", events.Delete())

		//TODO
		// router.Get("/edit", events.Edit())
		// router.Post("/edit", events.EditStore())
	})

	// Play
	handler.Route("/topics/{topicID}/play", func(router chi.Router) {
		router.Get("/1", play.Phase1())
		router.Get("/2", play.Phase2())
		router.Get("/3", play.Phase3())
		router.Post("/3", play.Store())
		router.Get("/review", play.Review())
	})

	// Scores
	handler.Route("/scores", func(router chi.Router) {
		router.Get("/", scores.List())
		router.Post("/", scores.Store())
	})

	// Users
	handler.Route("/users", func(router chi.Router) {
		router.Get("/register", users.Register())
		router.Post("/register", users.RegisterSubmit())
		router.Get("/login", users.Login())
		router.Post("/login", users.LoginSubmit())
		router.Get("/logout", users.Logout())
		router.Get("/{userID}/edit", users.EditUsername())
		router.Post("/{userID}", users.EditUsernameSubmit())
		router.Get("/{userID}/edit/password", users.EditPassword())
		router.Post("/{userID}", users.EditPasswordSubmit())

		//TODO
		// router.Get("/profile", users.Profile())
		// router.Get("/", users.List())
		// router.Post("/{userID}/delete", users.Delete())
	})

	return handler
}

// Handler
// Consists of the chi-multiplexer, a store interface and sessions
type Handler struct {
	*chi.Mux
	store    backend.Store
	sessions *scs.SessionManager
}

// Home
// A GET-method. Renders the home-page.
func (handler *Handler) Home() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData

		Topics []backend.Topic
		Scores []backend.Score
	}

	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/home.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute SQL statement to get topics
		tt, err := handler.store.Topics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get scores
		ss, err := handler.store.Scores(5, 0)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
			Topics:      tt,
			Scores:      ss,
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// About
// A GET-method. Renders the about-page.
func (handler *Handler) About() http.HandlerFunc {
	// Data to pass to HTML-templates
	type data struct {
		SessionData
	}
	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/about.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-templates with data
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(handler.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// withUser
// A middleware that replaces the potential user ID with a user object.
func (handler *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Retrieve user ID from session
		var userID int
		// User ID as interface
		userIDinf := handler.sessions.Get(req.Context(), "user_id")
		if userIDinf != nil {
			// If user ID interface isn't empty, turn it into an integer
			userID = userIDinf.(int)
		}

		// Execute SQL statement to get user
		user, err := handler.store.User(userID)
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
