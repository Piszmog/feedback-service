package transport_test

import (
	"errors"
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/model"
)

type mockDB struct {
	existsError bool
	exists      bool
	insertError bool
	findError   bool
	feedbacks   []model.Feedback
}

func (m mockDB) Exists(userID string, sessionID string) (bool, error) {
	if m.existsError {
		return false, errors.New("failed to check existence")
	}
	return m.exists, nil
}

func (m mockDB) Insert(feedback model.Feedback) error {
	if m.insertError {
		return errors.New("failed to insert")
	}
	return nil
}

func (m mockDB) Find(sessionID string, sort db.Sort, limit int) ([]model.Feedback, error) {
	if m.findError {
		return nil, errors.New("failed to find feedback")
	}
	return m.feedbacks, nil
}

func (m mockDB) FindWithFilter(sessionID string, filter db.Filter, sort db.Sort, limit int) ([]model.Feedback, error) {
	if m.findError {
		return nil, errors.New("failed to find feedback")
	}
	return m.feedbacks, nil
}

func (m mockDB) Close() {}
