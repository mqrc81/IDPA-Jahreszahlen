package web

/*
 * handler.go contains all HTTP-routes and is the basis for all HTTP-handlers
 */

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
	// FuncMap stores functions to use in HTML-template
	FuncMap = template.FuncMap{
		"rank": func(num int, page int, limit int) int {
			return (page-1)*limit + num + 1
		},
		"increment": func(num int) int {
			return num + 1
		},
	}
)

/*
 * NewHandler creates a new handler, including routes and middleware
 */
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

	h.Use(middleware.Logger)
	h.Use(sessions.LoadAndSave)
	h.Use(h.withUser)

	h.Get("/", h.Home())
	h.Get("/about", h.About())

	// Topics
	h.Route("/topics", func(r chi.Router) {
		r.Get("/", topics.List())
		r.Get("/new", topics.Create())
		r.Post("/", topics.Store())
		r.Post("/{topicID}/delete", topics.Delete())
		r.Get("/{topicID}/edit", topics.Edit())
		r.Get("/{topicID}", topics.Show())
	})

	// Events
	h.Route("/topics/{topicID}/events", func(r chi.Router) {
		r.Get("/new", events.Create())
		r.Post("/", events.Store())
		r.Post("/{eventID}/delete", events.Delete())
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
		r.Post("/logout", users.Logout())
		// r.Get("/profile", users.Profile())
		// r.Get("/", users.List())
		// r.Post("/{userID}/delete", users.Delete())
		// r.Get("/{userID}/edit", users.Edit())
	})

	return h
}

/*
 * Handler consists of the chi-multiplexer and a store
 */
type Handler struct {
	*chi.Mux
	store    backend.Store
	sessions *scs.SessionManager
}

/*
 * Home is a GET method that shows Homepage
 */
func (h *Handler) Home() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData

		Topics      []backend.Topic
		Scores      []backend.Score
		TopicsCount int
		EventsCount int
		UsersCount  int
	}

	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(homeHTML))

	return func(res http.ResponseWriter, req *http.Request) {
		// Execute SQL statement and return topics
		tt, err := h.store.Topics()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statement and return scores
		ss, err := h.store.Scores(5, 0)
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute SQL statements and return number of topics, events and users
		tCount, err := h.store.TopicsCount()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		eCount, err := h.store.EventsCount()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
		uCount, err := h.store.UsersCount()
		if err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}

		// Execute HTML-template
		if err := tmpl.Execute(res, data{
			SessionData: GetSessionData(h.sessions, req.Context()),
			Topics:      tt,
			Scores:      ss,
			TopicsCount: tCount, // TEMP
			EventsCount: eCount, // TEMP
			UsersCount:  uCount, // TEMP
		}); err != nil {
			http.Error(res, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * About is a GET method that shows information about this project
 */
func (h *Handler) About() http.HandlerFunc {
	// Data to pass to HTML-template
	type data struct {
		SessionData
	}
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(aboutHTML))

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

func (h *Handler) withUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		// Retrieve user ID from session
		var userID int
		userIDinf := h.sessions.Get(req.Context(), "user_id")
		if userIDinf != nil {
			userID = userIDinf.(int)
		}

		// Execute SQL statement
		user, err := h.store.User(userID)
		if err != nil {
			next.ServeHTTP(res, req)
			return
		}

		ctx := context.WithValue(req.Context(), "user", user)
		next.ServeHTTP(res, req.WithContext(ctx))
	})
}
