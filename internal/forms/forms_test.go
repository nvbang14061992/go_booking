package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	// create a new form with empty data
	newForm := Form{
		Errors: errors(map[string][]string{}),
	}

	if !newForm.Valid() {
		t.Error("got non-zero errors when there should be none")
	}
	
	// create a new form with httptest data, this is recommended way, 
	// because it simulates a real HTTP request and allows us to test the form handling in a more realistic way.
	r := httptest.NewRequest("POST", "/submit", nil)
	form := New(r.PostForm)
	if !form.Valid() {
		t.Error("got non-zero errors when there should be none")
	}
}

func TestForm_New(t *testing.T) {
	// create a new form with some data
	data := url.Values{}
	data.Add("name", "John")
	
	newForm := New(data)

	// check if the form's data is the same as the input data
	if newForm.Get("name") != "John" {
		t.Error("form data does not match input data")
	}
}

func TestForm_Required(t *testing.T) {
	// create a new form with some data
	data := url.Values{}
	data.Add("name", "John")
	data.Add("email", "")

	newForm := New(data)

	// check if the form is Required for "name" and "email"
	newForm.Required("name", "email")

	// the value of name is not empty, so there should be no error for "name"
	if newForm.Errors.Get("name") != "" {
		t.Error("got an error for name when there should be none")
	}
	// otherwise, the value of email is empty, so there should be an error for "email"
	if newForm.Errors.Get("email") == "" {
		t.Error("got no error for email when there should be one")
	}

}

func TestForm_Has(t *testing.T) {
	// create a new form with some data
	data := url.Values{}
	data.Add("name", "John")
	
	newForm := New(data)


	if !newForm.Has("name") {
		t.Error("form does not have 'name' field")
	}
	if newForm.Errors.Get("name") != "" {
		t.Error("got an error for name when there should be none")
	}
	if newForm.Has("email") {
		t.Error("form has 'email' field when it should not")
	}
	if newForm.Errors.Get("email") == "" {
		t.Error("got no error for email when there should be one")
	}

}

func TestForm_MinLength(t *testing.T) {
	// create a new form with some data
	data := url.Values{}
	data.Add("username", "john")
	data.Add("password", "123")
	
	// add data to form
	newForm := New(data)

	// check if the form has minimum length requirements
	newForm.MinLength("username", 3)
	newForm.MinLength("password", 5)

	// the value of username is not less than 3 characters, so there should be no error for "username"
	if newForm.Errors.Get("username") != "" {
		t.Error("got an error for username when there should be none")
	}
	// otherwise, the value of password is less than 5 characters, so there should be an error for "password"
	if newForm.Errors.Get("password") == "" {
		t.Error("got no error for password when there should be one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	// create a new form with some data
	data_valid := url.Values{}
	data_valid.Add("email", "john@example.com")
	
	// add data_valid to form
	newForm := New(data_valid)

	if !newForm.IsEmail("email") {
		t.Error("form does not have a valid email address")
	}
	if newForm.Errors.Get("email") != "" {
		t.Error("got an error for email when there should be none")
	}

	// create a new form with some data
	data_invalid := url.Values{}
	data_invalid.Add("email", "johnexample.com")
	// add data_invalid to form
	newForm_invalid := New(data_invalid)

	
	if newForm_invalid.IsEmail("email") {
		t.Error("form has an invalid email address")
	}
	if newForm_invalid.Errors.Get("email") == "" {
		t.Error("got no error for email when there should be one")
	}
}