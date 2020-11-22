package web

/*
 * TODO Header
 */

import (
	"html/template"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

var (
	funcMap = template.FuncMap{
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
	users := UserHandler{store: store, sessions: sessions}

	h.Use(middleware.Logger)
	h.Use(sessions.LoadAndSave)

	h.Get("/", h.Home())

	h.Route("/topics", func(r chi.Router) {
		r.Get("/", topics.List())
		r.Get("/new", topics.Create())
		r.Post("/store", topics.Store())
		r.Post("/{topicID}/delete", topics.Delete())
		r.Get("/{topicID}/edit", topics.Edit())
		r.Get("/{topicID}", topics.Show())
		r.Get("/{topicID}/play", topics.Play())
		r.Get("/{topicID}/scoreboard", topics.Scoreboard())

		r.Get("/{topicID}/events/new", events.Create())
		r.Post("/{topicID}/events/store", events.Store())
		r.Post("/{topicID}/events/{eventID}/delete", events.Delete())
	})
	h.Route("/users", func(r chi.Router) {
		r.Get("/register", users.Register())
		r.Post("/register", users.RegisterSubmit())
		//r.Get("/login", users.Login())
		//r.Get("/{username}", users.Profile())
		//r.Get("/{username}/edit", users.Edit())
		//r.Get("/{username}/scoreboard", users.Scoreboard())
		//r.Post("/store", users.Store())
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
 * Home TODO
 */
func (h *Handler) Home() http.HandlerFunc {
	// Parse HTML-template
	tmpl := template.Must(template.New("").Parse(homeHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		// Execute HTML-template
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}
