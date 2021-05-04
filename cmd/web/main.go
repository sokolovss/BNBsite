package main

import (
	"fmt"
	"github.com/sokolovss/BNBsite/pkg/config"
	"github.com/sokolovss/BNBsite/pkg/handlers"
	"github.com/sokolovss/BNBsite/pkg/render"
	"log"
	"net/http"
)

const portN = ":8080"

func main() {
	var app config.AppConfig

	tc, err := render.NewTemplateCache()
	if err != nil {
		log.Fatal("Cannot create templates cache")
	}

	app.TemplateCache = tc
	app.UseCache = false

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
