package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sokolovss/BNBsite/pkg/config"
	"github.com/sokolovss/BNBsite/pkg/handlers"
	"net/http"
)

func routes(app *config.AppConfig) http.Handler {
	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/contacts", handlers.Repo.Contacts)
	mux.Get("/generals", handlers.Repo.Generals)
	mux.Get("/colonels", handlers.Repo.Colonels)
	mux.Get("/reservation", handlers.Repo.Reservation)

	mux.Get("/search-availability", handlers.Repo.SearchAvailability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Get("/search-availability-json", handlers.Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}
