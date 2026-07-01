package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func Setup() http.Handler {
	r := chi.NewRouter()

	// ==========================
	// Global Middleware
	// ==========================

	// ==========================
	// Public Routes
	// ==========================

	// ==========================
	// Protected Routes
	// ==========================

	// ==========================
	// Admin Routes
	// ==========================

	return r
}
