package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	myH := myHandler{}
	h := NoSurf(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf("NoSurf returned wrong type: %T", v)
	}
}

func TestSessionLoad(t *testing.T) {
	myH := myHandler{}
	h := SessionLoad(&myH)

	switch v := h.(type) {
	case http.Handler:
		// do nothing, test passed
	default:
		t.Errorf("SessionLoad returned wrong type: %T", v)
	}
}