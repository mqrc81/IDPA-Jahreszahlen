package web

import (
	"html/template"
	"log"
	"net/http"
	"strconv"
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
		Mux: chi.NewMux(),
		store: store,
	}

	h.Use(middleware.Logger)
	h.Route("/units", func(router chi.Router) {
		router.Get("/", h.UnitsList())
		router.Get("/new", h.UnitsCreate())
	})

	return h
}

type Handler struct {
	*chi.Mux
	store backend.Store
}

// Temporary template for 'unitsList()'
const unitsListHTML = `
<h1>Themen</h1>
<dl>
    {{range .Units}}
        <dt><strong>{{.Title}} ({{.YearStart}} - {{.YearEnd}})</strong></dt>
        <dd>{{.Description}}</dd>
        <dd>Times played: {{.Playcount}}</dd>
    {{end}}
</dl>
`

// Temporary template for 'unitsCreate()'
const unitsCreateHTML = `
<h1>Neues Thema</h1>
<form action="/units" method="POST">
    <table>
        <tr>
            <td>Titel</td>
            <td><input type="text" name="title"/></td>
        </tr>
        <tr>
            <td>Zeitspanne</td>
            <td><input type="number" name="start_year"/> - <input type="number" name="end_year"></td>
        </tr>
        <tr>
            <td>Beschreibung (optional)</td>
            <td><input type="text" name="description"/></td>
        </tr>
    </table>
    <button type="submit">Thema erstellen</button>
</form>
`

/*
 * List of all units according to 'units_list.html' template
 */
func (h *Handler) UnitsList() http.HandlerFunc {
	type data struct {
		Units []backend.Unit
	}

	return func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.New("").Parse(unitsListHTML)) //TODO: ParseFiles("_.html")
		uu, err := h.store.Units()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(w, data{Units: uu}); err != nil {
			log.Fatal(err)
		}
	}
}

/*
 * Form to create new unit according to 'units_create.html' template
 */
func (h *Handler) UnitsCreate() http.HandlerFunc {
	tmpl := template.Must(template.New("").Parse(unitsCreateHTML)) //TODO: ParseFiles("_.html")
	return func(w http.ResponseWriter, r *http.Request) {
		if err := tmpl.Execute(w, nil); err != nil {
			log.Fatal(err)
		}
	}
}

/*
 * Stores unit created in 'UnitsCreate()'
 */
func (h *Handler) UnitsStore() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		title := r.FormValue("title")
		startYear, _ := strconv.Atoi(r.FormValue("start_year"))
		endYear, _ := strconv.Atoi(r.FormValue("end_year"))
		description := r.FormValue("description")

		if err := h.store.CreateUnit(&backend.Unit{
			Title: title,
			StartYear: startYear,
			EndYear: endYear,
			Description: description,
		}); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}