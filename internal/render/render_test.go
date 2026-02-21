package render

import (
	"net/http"
	"testing"

	"github.com/bangn/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	rq, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(rq.Context(), "flash", "123")
	result := AddDefaultData(&td, rq)

	if result.Flash != "123" {
		t.Error("flash value of 123 not found in session")
	}
	
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "http://testing", nil)
	if err != nil {
		return nil, err
	}

	// init fake context data
	ctx := r.Context()
	// add fake session data to the request context, so that we can test the AddDefaultData function
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	// add context back to the request, so that we can test the AddDefaultData function
	r = r.WithContext(ctx)

	return r, nil
}