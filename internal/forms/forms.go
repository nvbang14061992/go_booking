package forms

import (
	"net/http"
	"net/url"
)

type Form struct {
	url.Values
	Errors errors
}

// New initializes a Form struct.
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// Form checks if a field is present in the form data.
func (f *Form) Has(field string, r *http.Request) bool {
	x := r.Form.Get(field)
	if x == "" {
		return false
	}

	return true
}