package main

import (
	"net/http"

	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/handlers"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func routes(app *config.AppConfig) http.Handler {
	// mux := pat.New()
	mux := chi.NewRouter()
	
	// Use Chi middlewares
	// recover from panics (user request ), log request info, set secure headers
	// It handles the error gracefully by returning an http.StatusInternalServerError (500)
	// a bit similar to what next(err) function does in nodejs/express, where all the subsequent handlers/middlewares are skipped
	mux.Use(middleware.Recoverer)

	// protect against CSRF attacks
	mux.Use(NoSurf)

	// load and save session on every request
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/generals-quarters", handlers.Repo.Generals)
	mux.Get("/majors-suite", handlers.Repo.Majors)
	mux.Get("/search-availability", handlers.Repo.Availability)
	mux.Post("/search-availability", handlers.Repo.PostAvailability)
	mux.Post("/search-availability-json", handlers.Repo.AvailabilityJSON)
	mux.Get("/make-reservation", handlers.Repo.Reservation)
	mux.Post("/make-reservation", handlers.Repo.PostReservation)
	mux.Get("/reservation-summary", handlers.Repo.ReservationSummary)
	mux.Get("/contact", handlers.Repo.Contact)
	
	
	fileServer := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}