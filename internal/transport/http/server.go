package http

import (
	"net/http"

	"github.com/razedwell/go-hand/internal/transport/http/middleware"
)

type RouteRegistrar interface {
	RegisterRoutes(mux *http.ServeMux)
}

func NewServer(addr string, registrar ...RouteRegistrar) *http.Server {
	mux := http.NewServeMux()

	for _, r := range registrar {
		r.RegisterRoutes(mux)
	}

	handler := middleware.LoggingMiddleware(mux)

	return &http.Server{
		Addr:    addr,
		Handler: handler,
	}
}
