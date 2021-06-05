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
