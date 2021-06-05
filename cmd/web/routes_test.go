package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/sokolovss/BNBsite/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	h := routes(&app)

	switch h.(type) {
	case *chi.Mux:
	//pass
	default:
		t.Error("Type is not *chi.Mux")
	}
}
