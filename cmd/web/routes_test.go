package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/ldepner/bookings/internal/config"
	"testing"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig
	mux := routes(&app)

	switch mux.(type) {
	case *chi.Mux:
		//do nothing
	default:
		t.Error("type is not *chi.Mux")
	}
}
