package render

import (
	"encoding/gob"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/bangn/bookings/internal/config"
	"github.com/bangn/bookings/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	gob.Register(models.Reservation{})

	testApp.InProduction = false
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	testApp.Session = session
	app = &testApp

	os.Exit(m.Run())
}


// create fake http response writer
type myWriter struct {}

// if you look at the http.ResponseWriter interface, 
// it has three methods: Header(), Write([]byte) and WriteHeader(int), 
// so we need to implement these three methods for our myWriter struct, 
// so that we can use it as a fake response writer in our tests
// tips: just look at the http.ResponseWriter interface and implement the methods, 
// you don't need to care about the implementation details, 
// just return some dummy, empty data, 
// because we just want to test the AddDefaultData function, 
// we don't care about the response writer in our tests
func (w *myWriter) Header() http.Header {
	// because the Header() method returns a http.Header, which is a map[string][]string,
	// we can just return an empty map, because we don't care about the headers in our tests, 
	// we just need to implement the method to satisfy the http.ResponseWriter interface
	return http.Header{}
}

func (w *myWriter) Write(b []byte) (int, error) {
	return len(b), nil
}

func (w *myWriter) WriteHeader(statusCode int) {
}