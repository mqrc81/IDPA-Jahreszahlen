package web

import (
	"github.com/go-chi/chi/middleware"
	"html/template"
	"log"
	"net/http"

	"github.com/go-chi/chi"

	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
)

func NewHandler(store backend.Store) *Handler {
	handler := &Handler{
		Mux: chi.NewMux(),
		store: store,
	}

	handler.Use(middleware.Logger)
	handler.Route("/units", func(router chi.Router) {
		router.Get("/", handler.UnitsList())
	})

	return handler
}

type Handler struct {
	*chi.Mux
	store backend.Store
}

// Temporary
const unitsListHTML = `
<h1>Units</h1>
<dl>
{{range .}}
<dt><strong>{{.Title}}</strong> ({{.ID}})</dt>
<dd>{{.Description}}</dd>
<dd>Times played: {{.Playcount}}</dd>
{{end}}
</dl>
`

func (h *Handler) UnitsList() http.HandlerFunc {
	type data struct {
		Units []backend.Unit
	}

	return func(writer http.ResponseWriter, request *http.Request) {
		tmpl := template.Must(template.New("").Parse(unitsListHTML))
		uu, err := h.store.Units()
		if err != nil {
			http.Error(writer, err.Error(), http.StatusInternalServerError)
			return
		}

		if err := tmpl.Execute(writer, data{Units: uu}); err != nil {
			log.Fatal(err)
		}
	}
}