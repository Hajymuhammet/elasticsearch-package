package errors

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrIndexNotFound    = errors.New("elasticsearch: index not found")
	ErrDocumentNotFound = errors.New("elasticsearch: document not found")
	ErrInvalidQuery     = errors.New("elasticsearch: invalid query")
	ErrBulkFailed       = errors.New("elasticsearch: bulk operation failed")
	ErrInvalidID        = fmt.Errorf("invalid id provided")
)

// ElasticsearchError wraps ES error response with additional context.
type ElasticsearchError struct {
	Status  int    `json:"status"`
	Reason  string `json:"reason"`
	Details any    `json:"details,omitempty"`
}

func (e *ElasticsearchError) Error() string {
	return fmt.Sprintf("elasticsearch: status %d, reason %s", e.Status, e.Reason)
}

// Unwrap allows errors.Is / errors.As to work properly
func (e *ElasticsearchError) Unwrap() error {
	return fmt.Errorf("status %d: %s", e.Status, e.Reason)
}

// MarshalJSON makes error loggable in structured form
func (e *ElasticsearchError) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Type    string `json:"type"`
		Status  int    `json:"status"`
		Reason  string `json:"reason"`
		Details any    `json:"details,omitempty"`
	}{
		Type:    "ElasticsearchError",
		Status:  e.Status,
		Reason:  e.Reason,
		Details: e.Details,
	})
}

// WrapperElasticsearchError creates a new ElasticsearchError.
func WrapperElasticsearchError(status int, reason string, details any) *ElasticsearchError {
	return &ElasticsearchError{
		Status:  status,
		Reason:  reason,
		Details: details,
	}
}

// IsElasticsearchError checks whether an error is an ElasticsearchError.
func IsElasticsearchError(err error) (*ElasticsearchError, bool) {
	var esErr *ElasticsearchError
	if errors.As(err, &esErr) {
		return esErr, true
	}
	return nil, false
}
