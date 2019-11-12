package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Piszmog/feedback-service/model"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"time"
)

type MySQL struct {
	db *sql.DB
}

// Open create a connection to the MySQL DB.
func Open(options Options) (*MySQL, error) {
	if err := options.Validate(); err != nil {
		return nil, err
	}
	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=%t", options.Username, options.Password,
		options.Host, options.Port, options.DatabaseName, options.ParseTime))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to the MySQL DB: %w", err)
	}
	return &MySQL{db: db}, nil
}

// CreateFeedbackTableIfNotExists creates the 'feedback' table if it does not exist.
func (d MySQL) CreateFeedbackTableIfNotExists() error {
	_, err := d.db.Exec("CREATE TABLE IF NOT EXISTS `feedback`(" +
		"`id` INT UNSIGNED  NOT NULL AUTO_INCREMENT, " +
		"`userID` VARCHAR(255) NOT NULL, " +
		"`sessionID` VARCHAR(255) NOT NULL, " +
		"`comment` VARCHAR(255), " +
		"`rating` TINYINT NOT NULL, " +
		"`date` TIMESTAMP NOT NULL, " +
		"PRIMARY KEY (`id`), " +
		"INDEX(`userID`, `sessionID`), " +
		"INDEX(`sessionID`))")
	if err != nil {
		return fmt.Errorf("failed to create table 'feedback': %w", err)
	}
	return nil
}

func (d MySQL) Exists(userID string, sessionID string) (bool, error) {
	row := d.db.QueryRow("SELECT EXISTS(SELECT * FROM feedback WHERE userID=? AND sessionID=?)", userID, sessionID)
	var exists bool
	if err := row.Scan(&exists); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}
	return exists, nil
}

func (d MySQL) Insert(feedback model.Feedback) error {
	_, err := d.db.Exec("INSERT INTO feedback(`userID`, `sessionID`, `comment`, `rating`, `date`) VALUES (?,?,?,?,?)",
		feedback.UserID, feedback.SessionID, feedback.Comment, feedback.Rating, time.Now())
	if err != nil {
		return err
	}
	return nil
}

func (d MySQL) Find(sessionID string, sort Sort, limit int) ([]model.Feedback, error) {
	query := fmt.Sprintf("SELECT * FROM feedback where sessionID=? ORDER BY `date` %s LIMIT %d", sort, limit)
	rows, err := d.db.Query(query, sessionID)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	var feedback []model.Feedback
	for rows.Next() {
		var row model.Feedback
		if err := rows.Scan(&row.ID, &row.UserID, &row.SessionID, &row.Comment, &row.Rating, &row.Date); err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}
		feedback = append(feedback, row)
	}
	return feedback, nil
}

func (d MySQL) FindWithFilter(sessionID string, filter Filter, limit int, sort Sort) ([]model.Feedback, error) {
	query := fmt.Sprintf("SELECT * FROM feedback where sessionID=? AND rating =? ORDER BY `date` %s LIMIT %d", sort, limit)
	rows, err := d.db.Query(query, sessionID, filter.Rating)
	if err != nil {
		return nil, err
	}
	defer closeRows(rows)
	var feedback []model.Feedback
	for rows.Next() {
		var row model.Feedback
		if err := rows.Scan(&row.ID, &row.UserID, &row.SessionID, &row.Comment, &row.Rating, &row.Date); err != nil {
			return nil, fmt.Errorf("failed to read row: %w", err)
		}
		feedback = append(feedback, row)
	}
	return feedback, nil
}

func (d *MySQL) Close() {
	if err := d.db.Close(); err != nil {
		log.Println(fmt.Errorf("failed to close the connect to the DB: %w", err))
	}
}

func closeRows(rows *sql.Rows) {
	if err := rows.Close(); err != nil {
		log.Println(fmt.Errorf("failed to close MySQL rows: %w", err))
	}
}
