package db_test

import (
	"github.com/Piszmog/feedback-service/db"
	"testing"
)

var optionsTable = []struct {
	pass    bool
	options db.Options
}{
	{
		pass: true,
		options: db.Options{
			Username:     "user",
			Password:     "pass",
			Host:         "localhost",
			Port:         "8080",
			DatabaseName: "test",
			ParseTime:    false,
		},
	},
	{
		pass: false,
		options: db.Options{
			Password:     "pass",
			Host:         "localhost",
			Port:         "8080",
			DatabaseName: "test",
			ParseTime:    false,
		},
	},
	{
		pass: false,
		options: db.Options{
			Username:     "user",
			Host:         "localhost",
			Port:         "8080",
			DatabaseName: "test",
			ParseTime:    false,
		},
	},
	{
		pass: false,
		options: db.Options{
			Username:     "user",
			Password:     "pass",
			Port:         "8080",
			DatabaseName: "test",
			ParseTime:    false,
		},
	},
	{
		pass: false,
		options: db.Options{
			Username:     "user",
			Password:     "pass",
			Host:         "localhost",
			DatabaseName: "test",
			ParseTime:    false,
		},
	},
	{
		pass: false,
		options: db.Options{
			Username:  "user",
			Password:  "pass",
			Host:      "localhost",
			Port:      "8080",
			ParseTime: false,
		},
	},
}

func TestOptions_Validate(t *testing.T) {
	for _, entry := range optionsTable {
		err := entry.options.Validate()
		if entry.pass && err != nil {
			t.Errorf("expected options %+v to pass", entry.options)
		} else if !entry.pass && err == nil {
			t.Errorf("expected options %+v to not pass", entry.options)
		}
	}
}
