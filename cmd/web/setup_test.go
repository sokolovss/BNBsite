package main

import (
	"net/http"
	"os"
	"testing"
)

func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type TestHandler struct{}

func (t *TestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}
