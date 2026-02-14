package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/handlers"
	"github.com/bangn/bookings/internal/models"
	"github.com/bangn/bookings/internal/render"
)

const portNumber = ":8080"
var app config.AppConfig
var session *scs.SessionManager


// main is the main function of the application
func main() {
	err := run()
	
	if err != nil {
		log.Fatal(err)
	}

	// ---------------------------------------------
	// set up routes
	// ---------------------------------------------
	fmt.Printf("Starting server at port %s\n", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}

	// listen for a request
	err = srv.ListenAndServe()
	log.Fatal(err)
}

func run() error {
	// ---------------------------------------------
	// add place to store objects in session
	// ---------------------------------------------
	gob.Register(models.Reservation{})
	// This is necessary because the session manager needs to know how to encode and decode the Reservation struct when storing and retrieving it from the session.
	// Remember that, we only put map of any(map[string]interface{}) in template data, thus we need to register the struct that we want to put in the map, 
	// so that session can serialize and deserialize it correctly.
	// By registering the Reservation struct with gob, we ensure that the session manager can handle it properly when we store it in the session and retrieve it later.
	// If we don't register the struct, we may encounter errors when trying to store or retrieve Reservation objects from the session, because Go does not have built-in serialization for custom types, so it does not understand how to convert the Reservation struct to a format suitable for storage in the session (like a byte slice) and back.
	// consequently, we register any custom types that we plan to store in the session to avoid serialization issues.
	// ---------------------------------------------
	// The encoding/gob package is used for this purpose, and by registering the struct, we ensure that the session manager can handle it correctly.


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
		return err
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


	return nil
}