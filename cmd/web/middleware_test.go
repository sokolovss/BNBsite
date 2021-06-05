package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var testH TestHandler

	h := NoSurf(&testH)

	switch h.(type) {
	case http.Handler:
	//pass
	default:
		t.Error("Type is not http.Handler")
	}
}

func TestSessionLoad(t *testing.T) {
	var testH TestHandler

	h := SessionLoad(&testH)

	switch h.(type) {
	case http.Handler:
	//pass
	default:
		t.Error("Type is not http.Handler")
	}
}
