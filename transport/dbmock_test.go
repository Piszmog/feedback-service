package transport_test

import (
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/model"
)

type MockDB struct {
}

func (m MockDB) Exists(userID string, sessionID string) (bool, error) {
	panic("implement me")
}

func (m MockDB) Insert(feedback model.Feedback) error {
	panic("implement me")
}

func (m MockDB) Find(sessionID string, sort db.Sort, limit int) ([]model.Feedback, error) {
	panic("implement me")
}

func (m MockDB) FindWithFilter(sessionID string, filter db.Filter, sort db.Sort, limit int) ([]model.Feedback, error) {
	panic("implement me")
}

func (m MockDB) Close() {
	panic("implement me")
}
