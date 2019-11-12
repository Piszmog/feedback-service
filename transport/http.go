package transport

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/model"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"time"
)

const (
	contentTypeJSON   = "application/json"
	findLimit         = 15
	headerContentType = "Content-Type"
	headerUserID      = "Ubi-UserId"
	pathSessionID     = "sessionID"
	queryRating       = "rating"
)

// HTTPServer is the HTTP server with the provided configurations.
type HTTPServer struct {
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
	//
	// Add logging middleware to get some visibility into requests being made
	//
	router.Use(loggingMiddleware)
	//
	// Setup the possible paths
	//
	router.HandleFunc("/{sessionID}", s.insertFeedback()).Methods(http.MethodPost)
	router.HandleFunc("/{sessionID}", s.retrieveFeedback()).Methods(http.MethodGet)
	//
	// Configure the server
	//
	srv := &http.Server{
		Addr:         "0.0.0.0:" + s.Port, //todo address
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
		// Log the route that is being called
		//
		log.Println(r.RequestURI)
		//
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		//
		next.ServeHTTP(w, r)
	})
}

func (s *HTTPServer) insertFeedback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerContentType, contentTypeJSON)
		sessionID := mux.Vars(r)[pathSessionID]
		userID := r.Header.Get(headerUserID)
		//
		// Validate the user ID header
		//
		if len(r.Header.Get(headerUserID)) == 0 {
			writeHTTPError(http.StatusBadRequest, fmt.Sprintf("Missing Header '%s'", headerUserID), nil, w)
			return
		}
		//
		// Check if user has already submitted feedback for the session
		//
		exists, err := s.DB.Exists(userID, sessionID)
		if err != nil {
			writeHTTPError(http.StatusInternalServerError,
				fmt.Sprintf("Failed to check if user %s has previouly submitted feedback for session %s", userID, sessionID),
				err, w)
			return
		}
		if exists {
			writeHTTPError(http.StatusConflict,
				fmt.Sprintf("User %s has already submitted feedback for session %s", userID, sessionID), nil, w)
			return
		}
		//
		// Deserialize the request payload
		//
		defer closeRequestBody(r.Body)
		var feedback model.Feedback
		if err := json.NewDecoder(r.Body).Decode(&feedback); err != nil {
			writeHTTPError(http.StatusInternalServerError,
				fmt.Sprintf("Failed to decode user %s feedback for session %s", userID, sessionID), err, w)
			return
		}
		//
		// If user has not submitted feedback yet, insert their feedback
		//
		feedback.UserID = userID
		feedback.SessionID = sessionID
		feedback.Date = time.Now()
		if err := s.DB.Insert(feedback); err != nil {
			writeHTTPError(http.StatusInternalServerError,
				fmt.Sprintf("Failed to insert user %s feedback for session %s", userID, sessionID), err, w)
			return
		}
	}
}

func closeRequestBody(body io.ReadCloser) {
	if err := body.Close(); err != nil {
		log.Println(fmt.Errorf("failed to close the requeest body: %w", err))
	}
}

func (s *HTTPServer) retrieveFeedback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerContentType, contentTypeJSON)
		sessionID := mux.Vars(r)[pathSessionID]
		ratingFilter := r.URL.Query().Get(queryRating)
		var err error
		var feedback []model.Feedback
		//
		// If rating is provided in the query params, use it to find matching feedback.
		// Find the last 15 most recent feedback provided for the session.
		//
		if len(ratingFilter) > 0 {
			feedback, err = s.DB.FindWithFilter(sessionID, db.Filter{Rating: ratingFilter}, findLimit, db.Descending)
		} else {
			feedback, err = s.DB.Find(sessionID, db.Descending, findLimit)
		}
		if err != nil {
			writeHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to retrieve feedback for session %s", sessionID),
				err, w)
			return
		}
		//
		// Send data
		//
		if err := json.NewEncoder(w).Encode(feedback); err != nil {
			writeHTTPError(http.StatusInternalServerError, fmt.Sprintf("Failed to write feedback from session %s", sessionID),
				err, w)
			return
		}
	}
}

func writeHTTPError(statusCode int, reason string, err error, w http.ResponseWriter) {
	httpError := HTTPError{
		Code:   statusCode,
		Reason: reason,
		Err:    err,
	}
	w.WriteHeader(statusCode)
	if _, err := w.Write([]byte(httpError.ErrorJSON())); err != nil {
		log.Println(fmt.Errorf("failed to write HTTP error: %s: %w", httpError.Error(), err))
	}
	log.Println(httpError.Error())
	return
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
