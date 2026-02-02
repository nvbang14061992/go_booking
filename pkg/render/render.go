package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/bangn/bookings/pkg/config"
	"github.com/bangn/bookings/pkg/models"
	"github.com/justinas/nosurf"
)

// app is the application config variable,
// this will be set by the main application
// app here is global var in this module, not the same app in main
var app *config.AppConfig

// NewTemplates sets the config from the main application 
// for the render package
func NewTemplates(a *config.AppConfig) {
	// set the app config for the render package
	app = a
}

// AddDefaultData adds default data to all templates
func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate renders HTML templates
func RenderTemplate(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) {
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
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	myCache := map[string]*template.Template{}

	// get all of the files named *.page.tmpl from ./templates
	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}

	// range through all files ending with *.page.tmpl
	for _, page := range pages {
		name := filepath.Base(page) // get the file name
		// parse the page template file
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err // return added caches if exists and error
		}
		// look for layout templates  (*.layout.tmpl)
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		// if we found some, parse them
		if len(matches) > 0 {
			// parse the layout template file
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		// add the template to the cache map
		myCache[name] = ts
	}
	return myCache, nil // no error
}