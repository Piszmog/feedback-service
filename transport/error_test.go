package transport_test

import (
	"errors"
	"github.com/Piszmog/feedback-service/transport"
	"testing"
)

func TestHTTPError_Error(t *testing.T) {
	httpError := transport.HTTPError{
		Code:   404,
		Reason: "A Fail",
		Err:    errors.New("failed"),
	}
	msg := httpError.Error()
	if msg != "statusCode: 404, reason: A Fail: failed" {
		t.Errorf("error message does not match expected format: %s", msg)
	}
}

func TestHTTPError_ErrorJSON(t *testing.T) {
	httpError := transport.HTTPError{
		Code:   404,
		Reason: "A Fail",
		Err:    errors.New("failed"),
	}
	msg := httpError.ErrorJSON()
	if msg != `{"statusCode":404, "reason":"A Fail"}` {
		t.Errorf("error message does not match expected format: %s", msg)
	}
}
