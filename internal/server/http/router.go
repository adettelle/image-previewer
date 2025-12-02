package internalhttp

import (
	"github.com/go-chi/chi/v5"
)

func NewRouter(h *ImageHandler) chi.Router {
	r := chi.NewRouter()

	r.Get("/", mainPage)
	r.Get("/fill/{width}/{height}/*", h.preview)

	return r
}
