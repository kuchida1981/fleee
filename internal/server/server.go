package server

import (
	"fmt"
	"io/fs"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/kosuke/fleee/internal/handler"
)

// Server sets up the HTTP router and coordinates starting the server
type Server struct {
	router              *chi.Mux
	accountHandler      *handler.AccountHandler
	journalEntryHandler *handler.JournalEntryHandler
	port                string
	webFS               fs.FS
}

// NewServer creates a new Server instance
func NewServer(port string, accountHandler *handler.AccountHandler, journalEntryHandler *handler.JournalEntryHandler, webFS fs.FS) *Server {
	s := &Server{
		router:              chi.NewRouter(),
		accountHandler:      accountHandler,
		journalEntryHandler: journalEntryHandler,
		port:                port,
		webFS:               webFS,
	}
	s.setupRoutes()
	return s
}

func (s *Server) setupRoutes() {
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP) //nolint:staticcheck
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.Timeout(60 * time.Second))

	// API routes prefix
	s.router.Route("/api", func(r chi.Router) {
		r.Mount("/accounts", s.accountHandler.Routes())
		r.Mount("/journal-entries", s.journalEntryHandler.Routes())
	})

	// Serve static files and fallback to SPA index.html
	s.router.Handle("/*", spaHandler(s.webFS))
}

func spaHandler(webFS fs.FS) http.Handler {
	fsrv := http.FileServer(http.FS(webFS))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if len(path) > 0 && path[0] == '/' {
			path = path[1:]
		}
		if path == "" {
			path = "index.html"
		}

		// Check if file exists in the embedded filesystem
		_, err := fs.Stat(webFS, path)
		if err != nil {
			// File does not exist, rewrite request path to index.html for SPA
			r.URL.Path = "/index.html"
		}

		fsrv.ServeHTTP(w, r)
	})
}

// Start runs the HTTP server listening on the configured port
func (s *Server) Start() error {
	addr := fmt.Sprintf(":%s", s.port)
	fmt.Printf("Starting HTTP server on %s\n", addr)
	return http.ListenAndServe(addr, s.router)
}
