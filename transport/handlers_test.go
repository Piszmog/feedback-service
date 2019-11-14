package transport_test

import (
	"bytes"
	"encoding/json"
	"github.com/Piszmog/feedback-service/model"
	"github.com/Piszmog/feedback-service/transport"
	"github.com/gorilla/mux"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHTTPServer_InsertFeedback(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{exists: false}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test", "rating":4}`)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Ubi-UserId", "123")
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
}

func TestHTTPServer_InsertFeedback_TooHighRating(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{exists: false}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test", "rating":6}`)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Ubi-UserId", "123")
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	expected := `{"statusCode":400, "reason":"User 123 submitted rating 6 is not within the allowed range of 1-5 for session 987"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestHTTPServer_InsertFeedback_MalformedRequest(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{exists: false}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test" "rating":4}`)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Ubi-UserId", "123")
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	expected := `{"statusCode":400, "reason":"Failed to decode user 123 feedback for session 987"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestHTTPServer_InsertFeedback_ExistsError(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{existsError: true}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test", "rating":4}`)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Ubi-UserId", "123")
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	expected := `{"statusCode":500, "reason":"Failed to check if user 123 has previously submitted feedback for session 987"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestHTTPServer_InsertFeedback_FeedbackAlreadyExists(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{exists: true}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test", "rating":4}`)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Ubi-UserId", "123")
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusConflict {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusConflict)
	}
	expected := `{"statusCode":409, "reason":"User 123 has already submitted feedback for session 987"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestHTTPServer_InsertFeedback_MissingHeader(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{exists: false}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test", "rating":4}`)))
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusBadRequest)
	}
	expected := `{"statusCode":400, "reason":"Missing Header 'Ubi-UserId'"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestHTTPServer_InsertFeedback_InsertFailure(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{exists: false, insertError: true}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodPost, "/987", bytes.NewReader([]byte(`{"comment":"A Test", "rating":4}`)))
	if err != nil {
		t.Fatal(err)
	}
	request.Header.Set("Ubi-UserId", "123")
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.InsertFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	expected := `{"statusCode":500, "reason":"Failed to insert user 123 feedback for session 987"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}

func TestHTTPServer_RetrieveFeedback(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{feedbacks: []model.Feedback{{
		ID:        1,
		UserID:    "123",
		SessionID: "987",
		Comment:   "A Test",
		Rating:    4,
		Date:      time.Date(2019, 11, 12, 21, 00, 00, 00, time.UTC),
	}}}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodGet, "/987", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.RetrieveFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var feedbacks []model.Feedback
	if err := json.NewDecoder(recorder.Body).Decode(&feedbacks); err != nil {
		t.Fatal(err)
	}
	if len(feedbacks) != 1 {
		t.Errorf("expected feedbacks to be size 1 but got %d", len(feedbacks))
	} else if feedbacks[0].ID != 1 {
		t.Errorf("expected feedback ID to be 1 but got %d", feedbacks[0].ID)
	} else if feedbacks[0].UserID != "123" {
		t.Errorf("expected feedback ID to be 123 but got %s", feedbacks[0].UserID)
	} else if feedbacks[0].SessionID != "987" {
		t.Errorf("expected feedback ID to be 987 but got %s", feedbacks[0].SessionID)
	} else if feedbacks[0].Comment != "A Test" {
		t.Errorf("expected feedback ID to be 'A Test' but got %s", feedbacks[0].Comment)
	} else if feedbacks[0].Rating != 4 {
		t.Errorf("expected feedback rating to be 4 but got %d", feedbacks[0].Rating)
	}
}

func TestHTTPServer_RetrieveFeedback_WithFilter(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{feedbacks: []model.Feedback{{
		ID:        1,
		UserID:    "123",
		SessionID: "987",
		Comment:   "A Test",
		Rating:    4,
		Date:      time.Date(2019, 11, 12, 21, 00, 00, 00, time.UTC),
	}}}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodGet, "/987?rating=5", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.RetrieveFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}
	var feedbacks []model.Feedback
	if err := json.NewDecoder(recorder.Body).Decode(&feedbacks); err != nil {
		t.Fatal(err)
	}
	if len(feedbacks) != 1 {
		t.Errorf("expected feedbacks to be size 1 but got %d", len(feedbacks))
	} else if feedbacks[0].ID != 1 {
		t.Errorf("expected feedback ID to be 1 but got %d", feedbacks[0].ID)
	} else if feedbacks[0].UserID != "123" {
		t.Errorf("expected feedback ID to be 123 but got %s", feedbacks[0].UserID)
	} else if feedbacks[0].SessionID != "987" {
		t.Errorf("expected feedback ID to be 987 but got %s", feedbacks[0].SessionID)
	} else if feedbacks[0].Comment != "A Test" {
		t.Errorf("expected feedback ID to be 'A Test' but got %s", feedbacks[0].Comment)
	} else if feedbacks[0].Rating != 4 {
		t.Errorf("expected feedback rating to be 4 but got %d", feedbacks[0].Rating)
	}
}

func TestHTTPServer_RetrieveFeedback_FindError(t *testing.T) {
	//
	// Create server
	//
	server := transport.HTTPServer{DB: mockDB{findError: true}}
	//
	// Create Request, recorder, and handler
	//
	request, err := http.NewRequest(http.MethodGet, "/987", nil)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	router := mux.NewRouter()
	router.HandleFunc("/{sessionID}", server.RetrieveFeedback())
	//
	// Serve
	//
	router.ServeHTTP(recorder, request)
	//
	// Perform checks
	//
	if status := recorder.Code; status != http.StatusInternalServerError {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusInternalServerError)
	}
	expected := `{"statusCode":500, "reason":"Failed to retrieve feedback for session 987"}`
	if recorder.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v", recorder.Body.String(), expected)
	}
}
