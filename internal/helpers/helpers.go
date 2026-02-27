package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/bangn/bookings/internal/config"
)

var app *config.AppConfig

// NewHelpers sets the app config for the helpers package
func NewHelpers(a *config.AppConfig) {
	app = a
}

// ClientError is a helper function to send (http) client error messages
func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

// ServerError is a helper function to send server (internal app's) error messages
func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}