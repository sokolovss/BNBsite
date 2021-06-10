package handlers

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/justinas/nosurf"
	"github.com/sokolovss/BNBsite/internal/config"
	"github.com/sokolovss/BNBsite/internal/models"
	"github.com/sokolovss/BNBsite/internal/render"
	"html/template"
	"log"
	"net/http"
	"path/filepath"
	"time"
)

var app config.AppConfig
var session *scs.SessionManager

var pathToTemplates = "./../../templates"
var functions = template.FuncMap{}

func getRoutes() http.Handler {
	gob.Register(models.Reservation{})
	///////

	app.IsProduction = false
	app.UseCache = true

	session = scs.New()
	session.Lifetime = 12 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	tc, err := render.NewTemplateCache()
	if err != nil {
		log.Println("Cannot create templates cache")
	}

	app.TemplateCache = tc

	repo := NewRepo(&app)
	NewHandler(repo)
	render.NewTemplate(&app)

	mux := chi.NewMux()

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/contacts", Repo.Contacts)
	mux.Get("/generals", Repo.Generals)
	mux.Get("/colonels", Repo.Colonels)
	mux.Get("/reservation", Repo.Reservation)
	mux.Post("/reservation", Repo.PostReservation)
	mux.Get("/reservation-summary", Repo.ReservationSummary)

	mux.Get("/search-availability", Repo.SearchAvailability)
	mux.Post("/search-availability", Repo.PostAvailability)
	mux.Post("/search-availability-json", Repo.AvailabilityJSON)

	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

//NoSurf provides CSRF to POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.IsProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad loads and saves the sessions
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

//NewTestTemplateCache creates template cache as a map
func NewTemplateCache() (map[string]*template.Template, error) {
	pCache := make(map[string]*template.Template)
	p, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return pCache, err
	}

	for _, v := range p {
		n := filepath.Base(v)

		ts, err := template.New(n).Funcs(functions).ParseFiles(v)
		if err != nil {
			return pCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return pCache, err
		}

		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return pCache, err
			}
		}
		pCache[n] = ts
	}
	return pCache, nil
}
