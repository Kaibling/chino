package movies

import "github.com/go-chi/chi/v5"

func BuildRoute() *chi.Mux {
	r := chi.NewRouter()
	r.Post("/", create)
	// r.Get("/", readAllSaved)
	// r.Get("/all", readAll)
	// r.Get("/{id}", read)
	// r.Put("/{id}", update)
	r.Delete("/{name}", delete)

	return r
}
