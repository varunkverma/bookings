package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/varunkverma/bookings/pkg/config"
	"github.com/varunkverma/bookings/pkg/handlers"
)

// using chi routing package
func routes(app *config.AppConfig) http.Handler {
	// create a mutliplexer, which is a http handler
	mux := chi.NewRouter()

	// middlewares
	mux.Use(middleware.Recoverer) // Gracefully absorb panics and prints the stack trace
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)

	return mux
}
