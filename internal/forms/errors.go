package forms

type errors map[string][]string

// Add adds an error message for a given form field.
func (e errors) Add(field, message string) {
	// {"user": ["not_found", "deleted"], "password": ["too_short"]}
	e[field] = append(e[field], message)
}

// Get retrieves the first error message for a given form field.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}

	// avoid to return all error messages, it causes confusion in the UI, because BE errors can be very very long with stack traces etc.
	return es[0]
}