package model

import "time"

// Feedback is the feedback a user can provide for a session.
type Feedback struct {
	ID        int32
	UserID    string
	SessionID string
	Comment   string
	Rating    int8
	Date      time.Time
}
