package dashboard

import "github.com/go-chi/chi/v5"

func BuildRoute() *chi.Mux {
	r := chi.NewRouter()
	r.Get("/", show)

	return r
}
