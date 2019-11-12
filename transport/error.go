package transport

import "fmt"

// HTTPError is an error for a HTTP failure.
type HTTPError struct {
	Code   int
	Reason string
	Err    error
}

// Error provides a reason and a code in a JSON format.
func (e HTTPError) Error() string {
	reason := e.Reason
	//
	// If an error was provided, add to the reason
	//
	if e.Err != nil {
		reason = fmt.Errorf("%s: %w", e.Reason, e.Err).Error()
	}
	return fmt.Sprintf(`{"statusCode":%d, "reason":"%s"}`, e.Code, reason)
}
