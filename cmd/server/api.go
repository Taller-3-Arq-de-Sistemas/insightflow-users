package main

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/config"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/auth"
	"github.com/Taller-3-Arq-de-Sistemas/insightflow-users/internal/users"
)

type application struct {
	config       config.Config
	usersHandler *users.Handler
	authHandler  *auth.Handler
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Everything is OK"))
	})

	r.Route("/api/v1", func(r chi.Router) {
		r.Route("/users", func(r chi.Router) {
			r.With(app.authHandler.AuthMiddleware, app.authHandler.AdminMiddleware).Post("/", app.usersHandler.CreateUser)
			r.With(app.authHandler.AuthMiddleware).Get("/", app.usersHandler.ListUsers)
			r.Get("/{id}", app.usersHandler.FindUserById)
			r.With(app.authHandler.AuthMiddleware).Patch("/{id}", app.usersHandler.UpdateUser)
			r.With(app.authHandler.AuthMiddleware, app.authHandler.AdminMiddleware).Delete("/{id}", app.usersHandler.DeleteUser)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/login", app.authHandler.Login)
			r.Post("/register", app.authHandler.Register)
			r.With(app.authHandler.AuthMiddleware).Post("/logout", app.authHandler.Logout)
			r.Post("/validate", app.authHandler.ValidateToken)
		})
	})

	return r
}
