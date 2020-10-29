package web

import (
	"github.com/go-chi/chi"
	"github.com/mqrc81/IDPA-Jahreszahlen/backend"
	"net/http"
)

func NewHandler(store backend.Store) *Handler {
	handler := &Handler{
		mux: chi.NewMux(),
		store: store,
	}
	return handler
}

type Handler struct {
	mux *chi.Mux
	store backend.Store
}

func (h *Handler) UnitsList() http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {

	}
}