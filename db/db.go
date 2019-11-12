package db

import "github.com/Piszmog/feedback-service/model"

// DB is an interface for abstracting the interact with a database.
type DB interface {
	// Exists check whether the user has provided feedback for the specified session.
	Exists(userID string, sessionID string) (bool, error)

	// Insert inserts a feedback.
	Insert(feedback model.Feedback) error

	// Find finds feedback for a session. Limit specifies how many of the most recent feedback are returned.
	Find(sessionID string, limit int, sort Sort) ([]model.Feedback, error)

	// Find finds feedback for a session and with the provided filter. Limit specifies how many of the most recent feedback are returned.
	FindWithFilter(sessionID string, filter Filter, limit int, sort Sort) ([]model.Feedback, error)
}

// Filter is an additional filter that can be applied when querying for feedback.
type Filter struct {
	Rating string
}

// Sort determines whether to sort the feedback by newest or oldest.
type Sort int8

const (
	Ascending Sort = iota
	Descending
)
