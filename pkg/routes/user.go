package routes

import (
	"github.com/Adedunmol/zephyr/pkg/handlers"
	"github.com/go-chi/chi/v5"
)

func SetupUserRoutes(m *chi.Mux) {

	userRouter := chi.NewRouter()

	userRouter.Post("/register", handlers.CreateUserHandler)
	userRouter.Post("/login", handlers.LoginUserHandler)

	m.Mount("/users", userRouter)
}
