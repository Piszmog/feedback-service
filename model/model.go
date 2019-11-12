package model

import "time"

// Feedback is the feedback a user can provide for a session.
type Feedback struct {
	ID        int32     `json:"id"`
	UserID    string    `json:"userId"`
	SessionID string    `json:"sessionId"`
	Comment   string    `json:"comment"`
	Rating    int8      `json:"rating"`
	Date      time.Time `json:"date"`
}
