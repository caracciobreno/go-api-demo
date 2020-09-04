package http

import "github.com/gorilla/mux"

// API defines a HTTP API
type API interface {

	// RegisterRoutes register the routes of this API on the given Router
	RegisterRoutes(router *mux.Router)
}
