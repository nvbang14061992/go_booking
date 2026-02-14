package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	// setup tests, e.g. initialize the app, create template cache, etc.
	

	os.Exit(m.Run())
}

type myHandler struct{}

func (h *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// dummy handler, do nothing, just to test the middleware
}