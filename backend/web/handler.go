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
	h := &Handler{
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
	h.Use(middleware.Logger)
	h.Use(sessions.LoadAndSave)
	h.Use(h.withUser)

	// Home
	h.Get("/", h.Home())
	h.Get("/about", h.About())

	// Topics
	h.Route("/topics", func(r chi.Router) {
		r.Get("/", topics.List())
		r.Get("/new", topics.Create())
		r.Post("/", topics.CreateStore())
		r.Post("/{topicID}/delete", topics.Delete())
		r.Get("/{topicID}/edit", topics.Edit())
		r.Get("/{topicID}", topics.Show())
	})

	// Events
	h.Route("/topics/{topicID}/events", func(r chi.Router) {
		r.Get("/new", events.Create())
		r.Post("/", events.Store())
		r.Post("/{eventID}/delete", events.Delete())

		// TODO
		// r.Get("/edit", events.Edit())
		// r.Post("/edit", events.EditStore())
		// r.Get("/", events.List())
	})

	// Play
	h.Route("/topics/{topicID}/play", func(r chi.Router) {
		r.Get("/1", play.Phase1())
		r.Get("/2", play.Phase2())
		r.Get("/3", play.Phase3())
		r.Post("/3", play.Store())
		r.Get("/review", play.Review())
	})

	// Scores
	h.Route("/scores", func(r chi.Router) {
		r.Get("/", scores.List())
		r.Post("/", scores.Store())
	})

	// Users
	h.Route("/users", func(r chi.Router) {
		r.Get("/register", users.Register())
		r.Post("/register", users.RegisterSubmit())
		r.Get("/login", users.Login())
		r.Post("/login", users.LoginSubmit())
		r.Get("/logout", users.Logout())
		r.Get("/{userID}/edit/password", users.EditPassword())
		r.Post("/{userID}", users.EditPasswordStore())

		// TODO
		// r.Get("/profile", users.Profile())
		// r.Get("/", users.List())
		// r.Get("/{userID}/edit", users.EditUsername())
		// r.Post("/{userID}", users.EditUsernameStore())
		// r.Post("/{userID}/delete", users.Delete())
	})

	return h
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
func (h *Handler) Home() http.HandlerFunc {
	// Data to pass to HTML-template
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
		tt, err := h.store.Topics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement to get scores
		ss, err := h.store.Scores(5, 0)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
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
func (h *Handler) About() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
	}
	// Parse HTML-templates
	tmpl := template.Must(template.ParseFiles(
		"frontend/templates/layout.html",
		"frontend/templates/about.html"))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

// withUser
// A middleware that replaces the potential user ID with a user object.
func (h *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Retrieve user ID from session
		var userID int
		userIDinf := h.sessions.Get(req.Context(), "user_id") // user ID as interface{}
		if userIDinf != nil {
			userID = userIDinf.(int)
		}

		// Execute SQL statement to get user
		user, err := h.store.User(userID)
		if err != nil {
			next.ServeHTTP(res, req)
			return
		}

		// Add the user logged in to the session
		ctx := context.WithValue(req.Context(), "user", user)

		// Serve HTTP with response-writer and request
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
