package models

import "github.com/bangn/bookings/internal/forms"

// TemplateData holds data sent from handlers to templates
type TemplateData struct {
	// add any data you want to pass to the templates here
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float64
	Data      map[string]interface{}
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}