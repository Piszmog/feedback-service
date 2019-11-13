package transport

import (
	"encoding/json"
	"fmt"
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/model"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
	"strings"
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

// InsertFeedback inserts a user's feedback for a session. If a user has already submitted feedback, a 409 is returned.
func (s *HTTPServer) InsertFeedback() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(headerContentType, contentTypeJSON)
		sessionID := mux.Vars(r)[pathSessionID]
		userID := strings.TrimSpace(r.Header.Get(headerUserID))
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
				fmt.Sprintf("Failed to check if user %s has previously submitted feedback for session %s", userID, sessionID),
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
			writeHTTPError(http.StatusBadRequest,
				fmt.Sprintf("Failed to decode user %s feedback for session %s", userID, sessionID), err, w)
			return
		}
		if feedback.Rating > 5 {
			writeHTTPError(http.StatusBadRequest,
				fmt.Sprintf("User %s tried to submit a rating higher than the max rating value '5' for session %s. Submitted rating %d", userID, sessionID, feedback.Rating), nil, w)
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

// RetrieveFeedback retrieves the last 15 feedbacks for a specified session.
func (s *HTTPServer) RetrieveFeedback() func(w http.ResponseWriter, r *http.Request) {
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
			feedback, err = s.DB.FindWithFilter(sessionID, db.Filter{Rating: ratingFilter}, db.Descending, findLimit)
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
