package transport

import (
	"context"
	"fmt"
	"github.com/Piszmog/feedback-service/db"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

// HTTPServer is the HTTP server with the provided configurations.
type HTTPServer struct {
	Host         string
	Port         string
	WriteTimeout time.Duration
	ReadTimeout  time.Duration
	IdleTimeout  time.Duration
	DB           db.DB
	srv          *http.Server
}

// Start starts the HTTP server.
func (s *HTTPServer) Start() error {
	//
	// Setup the routing
	//
	router := mux.NewRouter()
	router.Use(loggingMiddleware)
	//
	// Setup the possible paths
	//
	router.HandleFunc("/{sessionID}", s.InsertFeedback()).Methods(http.MethodPost)
	router.HandleFunc("/{sessionID}", s.RetrieveFeedback()).Methods(http.MethodGet)
	//
	// Configure the server
	//
	srv := &http.Server{
		Addr:         s.Host + ":" + s.Port,
		WriteTimeout: s.WriteTimeout,
		ReadTimeout:  s.ReadTimeout,
		IdleTimeout:  s.IdleTimeout,
		Handler:      router,
	}
	s.srv = srv
	//
	// Start the server
	//
	if err := srv.ListenAndServe(); err != nil {
		return fmt.Errorf("failed to start server on port %s: %w", s.Port, err)
	}
	return nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//
		// Log how long the handler took
		//
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("Method: %s, URI: %s, Duration: %dmu\n", r.Method, r.RequestURI, time.Since(start).Microseconds())
	})
}

// Shutdown shutdowns the server with the provided timeout.
func (s *HTTPServer) Shutdown(timeout time.Duration) {
	//
	// Create a deadline
	//
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	//
	// Will wait for timeout if there are connections
	//
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Println(err)
	}
}
