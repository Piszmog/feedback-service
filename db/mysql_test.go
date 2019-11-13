package db_test

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Piszmog/feedback-service/db"
	"github.com/Piszmog/feedback-service/model"
	"testing"
	"time"
)

func TestMySQL_CreateFeedbackTableIfNotExists(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS `feedback`*").WillReturnResult(sqlmock.NewResult(1, 1))
	//
	// Run the test
	//
	createError := mySQL.CreateFeedbackTableIfNotExists()
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if createError != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	mock.ExpectClose()
}

func TestMySQL_CreateFeedbackTableIfNotExists_WithError(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectExec("CREATE TABLE IF NOT EXISTS `feedback`*").WillReturnError(errors.New("failed"))
	//
	// Run the test
	//
	createError := mySQL.CreateFeedbackTableIfNotExists()
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if createError == nil {
		t.Error("expected error to occurred")
	}
	mock.ExpectClose()
}

func TestMySQL_Exists(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("^SELECT EXISTS\\(SELECT \\* FROM feedback*").WithArgs("12345", "98765").
		WillReturnRows(sqlmock.NewRows([]string{"exists"}).AddRow(true))
	//
	// Run the test
	//
	exists, existsError := mySQL.Exists("12345", "98765")
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if !exists {
		t.Error("expected the row to exist")
	} else if existsError != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	mock.ExpectClose()
}

func TestMySQL_Exists_NoRow(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("^SELECT EXISTS\\(SELECT \\* FROM feedback*").WithArgs("12345", "98765").
		WillReturnError(sql.ErrNoRows)
	//
	// Run the test
	//
	exists, existsError := mySQL.Exists("12345", "98765")
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if exists {
		t.Error("expected the row to not exist")
	} else if existsError != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	mock.ExpectClose()
}

func TestMySQL_Exists_WithError(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("^SELECT EXISTS\\(SELECT \\* FROM feedback*").WithArgs("12345", "98765").
		WillReturnError(errors.New("failed"))
	//
	// Run the test
	//
	exists, existsError := mySQL.Exists("12345", "98765")
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if exists {
		t.Error("expected the row to not exist")
	} else if existsError == nil {
		t.Error("expected error to occurred")
	}
	mock.ExpectClose()
}

func TestMySQL_Insert(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectExec("INSERT INTO feedback*").WithArgs("123", "987", "A Test", 5, anyTime{}).
		WillReturnResult(sqlmock.NewResult(1, 1))
	//
	// Run the test
	//
	insertError := mySQL.Insert(model.Feedback{
		UserID:    "123",
		SessionID: "987",
		Comment:   "A Test",
		Rating:    5,
		Date:      time.Now(),
	})
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if insertError != nil {
		t.Errorf("unexpected error occurred: %v", err)
	}
	mock.ExpectClose()
}

func TestMySQL_Insert_WithError(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectExec("INSERT INTO feedback*").WithArgs("123", "987", "A Test", 5, anyTime{}).
		WillReturnError(errors.New("failed"))
	//
	// Run the test
	//
	insertError := mySQL.Insert(model.Feedback{
		UserID:    "123",
		SessionID: "987",
		Comment:   "A Test",
		Rating:    5,
		Date:      time.Now(),
	})
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if insertError == nil {
		t.Error("expected error to occurred")
	}
	mock.ExpectClose()
}

func TestMySQL_Find(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("SELECT \\* FROM feedback where sessionID=\\? ORDER BY `date` DESC LIMIT 1").
		WithArgs("987").WillReturnRows(sqlmock.NewRows([]string{"id", "userID", "sessionID", "comment", "rating", "date"}).
		AddRow(1, "123", "987", "A Test", 5, time.Now()))
	//
	// Run the test
	//
	feedbacks, findError := mySQL.Find("987", db.Descending, 1)
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if findError != nil {
		t.Errorf("unexpected error occurred: %v", err)
	} else if feedbacks == nil {
		t.Error("no feedback returned")
	}
	mock.ExpectClose()
}

func TestMySQL_Find_WithError(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("SELECT \\* FROM feedback where sessionID=\\? ORDER BY `date` DESC LIMIT 1").
		WithArgs("987").WillReturnError(errors.New("failed"))
	//
	// Run the test
	//
	feedbacks, findError := mySQL.Find("987", db.Descending, 1)
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if findError == nil {
		t.Errorf("expected error to occurred")
	} else if feedbacks != nil {
		t.Error("feedbacks returned")
	}
	mock.ExpectClose()
}

func TestMySQL_Find_BadColumns(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("SELECT \\* FROM feedback where sessionID=\\? ORDER BY `date` DESC LIMIT 1").
		WithArgs("987").WillReturnRows(sqlmock.NewRows([]string{"id", "userID", "sessionID", "comment"}).
		AddRow("1", 123, "987", "A Test"))
	//
	// Run the test
	//
	feedbacks, findError := mySQL.Find("987", db.Descending, 1)
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if findError == nil {
		t.Error("expected error to occurred")
	} else if feedbacks != nil {
		t.Error("feedbacks returned")
	}
	mock.ExpectClose()
}

func TestMySQL_FindWithFilter(t *testing.T) {
	//
	// Mock the SQL DB
	//
	mySQL, mock := createMockDB(t)
	defer mySQL.Close()
	//
	// Setup Mocks
	//
	mock.ExpectQuery("SELECT \\* FROM feedback where sessionID=\\? AND rating=\\? ORDER BY `date` DESC LIMIT 1").
		WithArgs("987", "5").WillReturnRows(sqlmock.NewRows([]string{"id", "userID", "sessionID", "comment", "rating", "date"}).
		AddRow(1, "123", "987", "A Test", 5, time.Now()))
	//
	// Run the test
	//
	feedbacks, findError := mySQL.FindWithFilter("987", db.Filter{Rating: "5"}, db.Descending, 1)
	//
	// Ensure expectations were met
	//
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations %v", err)
	} else if findError != nil {
		t.Errorf("unexpected error occurred: %v", err)
	} else if feedbacks == nil {
		t.Error("no feedback returned")
	}
	mock.ExpectClose()
}

func createMockDB(t *testing.T) (*db.MySQL, sqlmock.Sqlmock) {
	connection, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	return &db.MySQL{DB: connection}, mock
}

type anyTime struct{}

// Match satisfies sqlmock.Argument interface
func (a anyTime) Match(v driver.Value) bool {
	_, ok := v.(time.Time)
	return ok
}
