package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/handlers"
	"github.com/bangn/bookings/internal/render"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager


// main is the main function of the application
func main() {
	// <<<<<<<<<<<<<<<<<<<<<<<<<<<<<
	// set this to true in production
	app.InProduction = false

	// ---------------------------------------------
	// create session configuration parameters// ---------------------------------------------
	// ---------------------------------------------
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	// ---------------------------------------------
	// create cache for templates to render later
	// ---------------------------------------------
	tc, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("Can not creae template cache")
		return
	}

	// ---------------------------------------------
	// add template cache to app config
	// ---------------------------------------------
	app.TemplateCache = tc
	app.UseCache = false // development mode, the template will be reloaded on every request, 
	// not use the global cache, because the templates may change frequently, but in production mode, set it to true,
	// because no one changes the templates in the backend frequently

	
	// ---------------------------------------------
	// set the app config to the render package
	// ---------------------------------------------
	render.NewTemplates(&app)

	// ---------------------------------------------
	// set the app config to the handler package, to render templates
	// ---------------------------------------------
	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)
	
	// ---------------------------------------------
	// set up routes
	// ---------------------------------------------
	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	fmt.Printf("Starting server at port %s\n", portNumber)
	// listen for a request
	err = srv.ListenAndServe()
	log.Fatal(err)
}