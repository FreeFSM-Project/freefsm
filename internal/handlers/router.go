package handlers

import (
	"net/http"

	"github.com/MartialM1nd/freefsm/internal/config"
	"github.com/MartialM1nd/freefsm/internal/middleware"
	"github.com/MartialM1nd/freefsm/internal/services"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func New(db *pgxpool.Pool, sessions *services.SessionService, cfg *config.Config) http.Handler {
	r := chi.NewRouter()

	authMW := middleware.Auth(sessions)

	r.Group(func(r chi.Router) {
		r.Use(authMW)
		r.Get("/", handleDashboard)
		r.Post("/logout", func(w http.ResponseWriter, r *http.Request) {
			handleLogout(w, r, sessions)
		})
	})

	authHandler := NewAuthHandler(db, sessions)
	r.Get("/login", authHandler.ServeHTTP)
	r.Post("/login", authHandler.ServeHTTP)

	setupHandler := NewSetupHandler(db, sessions, cfg)
	r.Get("/setup", setupHandler.ServeHTTP)
	r.Post("/setup", setupHandler.ServeHTTP)

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(`{"status":"ok"}`))
	})

	return r
}
