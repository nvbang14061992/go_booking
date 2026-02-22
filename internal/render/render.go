package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// app is the application config variable,
// this will be set by the main application
// app here is global var in this module, not the same app in main
var app *config.AppConfig

var functions = template.FuncMap{}

var pathToTemplates = "./templates"

// NewTemplates sets the config from the main application 
// for the render package
func NewTemplates(a *config.AppConfig) {
	// set the app config for the render package
	app = a
}

// AddDefaultData adds default data to all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// go can not indentify specific request's context, thus we need to pass the request context to session, 
	// so that session can identify which request's session data to get 
	// go philosophy: explicit is better than implicit, you will add only the things you want to add, not all the things in session, 
	// thus we need to pop the data we want to add to template data, 
	// and add them to template data, then pass the template data to template, 
	// this is more explicit than just pass the whole session data to template, 
	// which may contain some sensitive data that we don't want to expose to template
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	
	
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders HTML templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {
	var tc map[string]*template.Template

	if app.UseCache {
		// get the template cache from the app config, global app variable
		tc = app.TemplateCache
	} else {
		// create a new template cache
		tc, _ = CreateTemplateCache()
	}

	// get requested template from cache
	t, ok := tc[tmpl]
	if !ok {
		log.Fatal("Could not get template from template cache")
		return errors.New("could not get template from template cache")
	}
	buffer := new(bytes.Buffer)

	// add default data to template data in all templates
	td = AddDefaultData(td, r)

	// execute the template, meaning apply the template to the data
	_ = t.Execute(buffer, td)

	// render the template
	_, err := buffer.WriteTo(w)
	if err != nil {
		log.Fatal(">>>>>>Error writing buffer to response writer:", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplates))
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page) // get the file name
		// parse the page template file
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err // return added caches if exists and error
		}
		// look for layout templates  (*.layout.tmpl)
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
		if err != nil {
			return myCache, err
		}
		// if we found some, parse them
		if len(matches) > 0 {
			// parse the layout template file
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplates))
			if err != nil {
				return myCache, err
			}
		}
		// add the template to the cache map
		myCache[name] = ts
	}
	return myCache, nil // no error
}