package main

import (
	"encoding/gob"
	"fmt"
	"github.com/alexedwards/scs/v2"
	config "github.com/sokolovss/BNBsite/internal/config"
	"github.com/sokolovss/BNBsite/internal/driver"
	handlers "github.com/sokolovss/BNBsite/internal/handlers"
	"github.com/sokolovss/BNBsite/internal/helpers"
	"github.com/sokolovss/BNBsite/internal/models"
	render "github.com/sokolovss/BNBsite/internal/render"
	"log"
	"net/http"
	"os"
	"time"
)

const portN = ":8080"

var app config.AppConfig
var session *scs.SessionManager
var infoLog *log.Logger
var errorLog *log.Logger

func main() {

	db, err := run()
	defer db.SQL.Close()

	if err != nil {
		log.Fatal(err)
	}

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

func run() (*driver.DB, error) {
	//Defines what will be stored in session (primitives are already built in)
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	///////

	app.IsProduction = false
	app.UseCache = false

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 12 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.IsProduction

	app.Session = session

	//connect to DB
	log.Println("Connecting to db...")
	db, err := driver.ConnectSQL("host=localhost port=5432 dbname=bnbsite user=sergey password=")
	if err != nil {
		log.Fatal("Cannot connect to database", err)
	}
	log.Println("Connected to database")

	tc, err := render.NewTemplateCache()
	if err != nil {
		log.Println("Cannot create templates cache")
		return nil, err
	}

	app.TemplateCache = tc

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandler(repo)
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}
