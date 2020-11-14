package web

import (
	"net/http"
	"strconv"
	"text/template"
	//
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	//
	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

/*
 * Creates new handler, including routers and middleware
 */
func NewHandler(store backend.Store) *Handler {
	h := &Handler{
		Mux:   chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)

	h.Route("/units", func(r chi.Router) {

		r.Get("/", h.UnitsList())
		r.Get("/new", h.UnitsCreate())
		r.Post("/", h.UnitsStore())
		r.Post("/{id}/delete", h.UnitsDelete())
		r.Get("/{id}/edit", h.UnitsEdit())
		r.Get("/{id}", h.UnitsShow())
		r.Get("/{id}/events/new", h.EventsCreate())
		r.Post("/", h.EventsStore())
	})

	return h
}

type Handler struct {
	*chi.Mux
	store backend.Store
}

/*
 * List of all units
 */
func (h *Handler) UnitsList() http.HandlerFunc {
	type data struct {
		Units []backend.Unit
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("").Parse(unitsListHTML))

		uu, err := h.store.Units()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data{Units: uu}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Form to create new unit
 */
func (h *Handler) UnitsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(unitsCreateHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Stores unit created
 */
func (h *Handler) UnitsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		startYear, _ := strconv.Atoi(r.FormValue("start_year"))
		endYear, _ := strconv.Atoi(r.FormValue("end_year"))
		description := r.FormValue("description")

		if err := h.store.CreateUnit(&backend.Unit{
			ID:			 0,
			Title:       title,
			StartYear:   startYear,
			EndYear:     endYear,
			Description: description,
			PlayCount:   0,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of units
		http.Redirect(w, r, "/units", http.StatusFound)
	}
}

/*
 * Deletes unit
 */
func (h *Handler) UnitsDelete() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		if err := h.store.DeleteUnit(id); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Redirect to list of units
		http.Redirect(w, r, "/units", http.StatusFound)
	}
}

/*
 * Shows Unit with options to play, see leaderboard, (edit if admin)
 */
func (h *Handler) UnitsShow() http.HandlerFunc {
	type data struct {
		Unit backend.Unit
	}

	tmpl := template.Must(template.New("").Parse(unitsShowHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, err := h.store.Unit(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data{Unit: u}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Edit unit and its events
 */
func (h *Handler) UnitsEdit() http.HandlerFunc {
	type data struct {
		Unit   backend.Unit
		Events []backend.Event
	}

	tmpl := template.Must(template.New("").Parse(unitsEditHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		id, _ := strconv.Atoi(chi.URLParam(r, "id"))
		u, err := h.store.Unit(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		ee, err := h.store.EventsByUnit(id)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data{Unit: u, Events: ee}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Form to create new event
 */
func (h *Handler) EventsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(eventsCreateHTML))

	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

/*
 * Stores event created
 */
func (h *Handler) EventsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		unitID, _ := strconv.Atoi(chi.URLParam(r, "id"))
		title := r.FormValue("title")
		year, _ := strconv.Atoi(r.FormValue("year"))

		if err := h.store.CreateEvent(&backend.Event{
			ID:     0,
			UnitID: unitID,
			Title:  title,
			Year:   year,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}
