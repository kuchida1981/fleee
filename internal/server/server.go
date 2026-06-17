package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kosuke/fleee/internal/handler"
)

// Server sets up the HTTP router and coordinates starting the server
type Server struct {
	router         *chi.Mux
	accountHandler *handler.AccountHandler
	port           string
}

// NewServer creates a new Server instance
func NewServer(port string, accountHandler *handler.AccountHandler) *Server {
	s := &Server{
		router:         chi.NewRouter(),
		accountHandler: accountHandler,
		port:           port,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// API routes prefix
	s.router.Route("/api", func(r chi.Router) {
		r.Mount("/accounts", s.accountHandler.Routes())
	})

	// Fallback route for API server health check
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.Write([]byte("fleee API Server is running"))
	})
}

// Start runs the HTTP server listening on the configured port
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	fmt.Printf("Starting HTTP server on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
