package app

import (
	"context"
	"net/http"
	"time"
)

// httpServer defines a HTTP server that is provided by the app
type httpServer struct {
	handler http.Handler
	server  *http.Server
}

// newHTTPServer creates a httpServer with a http.Handler
func newHTTPServer(handler http.Handler, addr string) *httpServer {
	s := &httpServer{
		handler: handler,
		server: &http.Server{
			ReadTimeout:  1 * time.Second,
			WriteTimeout: 10 * time.Second,
			Addr:         addr,
			Handler:      handler,
		},
	}

	return s
}

// Start listens to incoming requests and servers them
func (s *httpServer) Start(context.Context) error {
	return s.server.ListenAndServe()
}

func (s *httpServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}
