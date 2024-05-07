package routes

import "github.com/go-chi/chi/v5"

func SetupRoutes() *chi.Mux {
	m := chi.NewRouter()

	SetupUserRoutes(m)

	return m
}
