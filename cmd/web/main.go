package main

import (
	"fmt"
	"github.com/alexedwards/scs/v2"
	config "github.com/sokolovss/BNBsite/internal/config"
	handlers "github.com/sokolovss/BNBsite/internal/handlers"
	render "github.com/sokolovss/BNBsite/internal/render"
	"log"
	"net/http"
	"time"
)

const portN = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {

	app.IsProduction = false
	app.UseCache = false

	session = scs.New()
	session.Lifetime = 12 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	tc, err := render.NewTemplateCache()
	if err != nil {
		log.Println(err)
		log.Fatal("Cannot create templates cache")
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app)
	handlers.NewHandler(repo)
	render.NewTemplate(&app)

	fmt.Printf("Starting on port %v\n", portN)

	srv := &http.Server{
		Addr:    portN,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal("Error  - starting the server", err)
	}
}
