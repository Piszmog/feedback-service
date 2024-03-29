package db

import (
	"errors"
	"github.com/Piszmog/feedback-service/model"
)

// DB is an interface for abstracting the interact with a database.
type DB interface {
	// Exists check whether the user has provided feedback for the specified session.
	Exists(userID string, sessionID string) (bool, error)

	// Insert inserts a feedback.
	Insert(feedback model.Feedback) error

	// Find finds feedback for a session. Limit specifies how many of the most recent feedback are returned.
	Find(sessionID string, sort Sort, limit int) ([]model.Feedback, error)

	// Find finds feedback for a session and with the provided filter. Limit specifies how many of the most recent feedback are returned.
	FindWithFilter(sessionID string, filter Filter, sort Sort, limit int) ([]model.Feedback, error)

	// Close closes the DB connection.
	Close()
}

// Filter is an additional filter that can be applied when querying for feedback.
type Filter struct {
	Rating string
}

// Sort determines whether to sort the feedback by newest or oldest.
type Sort string

const (
	// Ascending sorts results in ascending order
	Ascending Sort = "ASC"
	// Descending sorts results in descending order
	Descending Sort = "DESC"
)

// Options are the connection options used to connect to a DB.
type Options struct {
	Username     string
	Password     string
	Host         string
	Port         string
	DatabaseName string
	ParseTime    bool
}

// Validate validates the provided options.
func (o Options) Validate() error {
	if len(o.Username) == 0 {
		return errors.New("require username to connect to the DB")
	} else if len(o.Password) == 0 {
		return errors.New("require password to connect to the DB")
	} else if len(o.Host) == 0 {
		return errors.New("require host to connect to the DB")
	} else if len(o.Port) == 0 {
		return errors.New("require port to connect to the DB")
	} else if len(o.DatabaseName) == 0 {
		return errors.New("require database name to connect to the DB")
	}
	return nil
}
