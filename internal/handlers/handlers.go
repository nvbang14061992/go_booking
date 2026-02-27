package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/forms"
	"github.com/bangn/bookings/internal/helpers"
	"github.com/bangn/bookings/internal/models"
	"github.com/bangn/bookings/internal/render"
)

// Repo the repository used by the handlers
var Repo *Repository

// Repository is the repository type
type Repository struct{
	App *config.AppConfig
}

// NewRepo creates a new repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers sets the repository for the handlers package
func NewHandlers(r *Repository) {
	Repo = r
}
// Home is the home page handler
func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {
	remoteIP := r.RemoteAddr
	m.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, r, "home.page.tmpl", &models.TemplateData{})
}

// About is the about page handler
func (m *Repository) About(w http.ResponseWriter, r *http.Request) {
	stringMap := make(map[string]string)
	stringMap["test"] = "Hello, again. This is the about page"

	// example of trying to use session
	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.RenderTemplate(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

// Reservation renders the make reservation and display form
func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {
	var emptyReservation models.Reservation
	data := make(map[string]interface{})
	data["reservation"] = emptyReservation

	// models.Reservation{} was added to gob, thus we can store it in session,
	// and we can also get it from session, but here we just initialize an empty reservation struct, then pass it to template,
	render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		// initialize an empty form, then pass data to it, then return data or errors back via post handler
		Form: forms.New(nil),
		Data: data,
	})
}


// PostReservation handles the reservation form submission
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
	}

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		// models.Reservation{} was added to gob, thus we can store it in session,
		// here session can serialize the reservation struct, gotten from form, and pass it back to template, 
		// so that the user doesn't have to re-enter the data they already entered, just correct the errors
		data["reservation"] = reservation
		render.RenderTemplate(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// save reservation in session, then next page will get it from session via redirect
	m.App.Session.Put(r.Context(), "reservation", reservation)
	// redirect to summary page
	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)
}

// Generals renders the room page
func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "generals.page.tmpl", &models.TemplateData{})
}


// Majors renders the room page
func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "majors.page.tmpl", &models.TemplateData{})
}


// Availability render search availability page
func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// PostAvailability render search availability page
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	w.Write([]byte(fmt.Sprintf("start date is %s and end date is %s", start, end)))
}

type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// AvailabilityJSON handle request for availability and send JSON response
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      false,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		helpers.ServerError(w, err)
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Contact render search Contact page
func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.tmpl", &models.TemplateData{})
}

// ReservationSummary render reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	// the type assertion .(models.Reservation) checks if the value stored in the session under the key "reservation" is of type models.Reservation.
	// If the assertion is successful, ok will be true, and reservation will hold the value with the correct type.
	// If the assertion fails (meaning the value is not of the expected type), ok will be false, and reservation will be the zero value for models.Reservation.
	if !ok {
		m.App.ErrorLog.Println("Cannot get reservation from session")
		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	// this is type assertion technique in Go, example:
	// 	var x interface{} = models.Reservation{}

	// r, ok := x.(models.Reservation)

	// ok == true
	// r contains the struct

	// free up the session, 
	// since we already got the reservation data from session, we can remove it from session, 
	// so that it won't take up space in session, and also it won't cause confusion if we have multiple reservations in session, 
	// we only want to keep the current reservation in session, thus we remove it after we get it
	m.App.Session.Remove(r.Context(), "reservation")

	// create a map to hold data sent to template
	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.RenderTemplate(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data: data,
	})
}