package routes

import "github.com/go-chi/chi/v5"

func SetupRoutes() {
	m := chi.NewRouter()

	SetupUserRoutes(m)
}
