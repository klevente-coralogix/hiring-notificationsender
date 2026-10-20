// Package thirdparty contains the third-party API that messages are sent to,
// plus a mock for tests.
package thirdparty

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"
)

// Errors returned by the third-party API. They mirror its HTTP error responses.
var (
	// ErrBadRequest (400): the message is malformed.
	ErrBadRequest = errors.New("400 Bad Request")
	// ErrUnauthorized (401): our credentials are missing or invalid.
	ErrUnauthorized = errors.New("401 Unauthorized")
	// ErrPayloadTooLarge (413): the message is too large.
	ErrPayloadTooLarge = errors.New("413 Payload Too Large")
	// ErrInternalServerError (500): the API failed while processing the message.
	ErrInternalServerError = errors.New("500 Internal Server Error")
	// ErrServiceUnavailable (503): the API is temporarily down or overloaded.
	ErrServiceUnavailable = errors.New("503 Service Unavailable")
	// ErrGatewayTimeout (504): an upstream of the API timed out. The message may
	// or may not have been delivered.
	ErrGatewayTimeout = errors.New("504 Gateway Timeout")
)

// TooManyRequestsError (429): we are sending too fast. Use errors.As to get it.
type TooManyRequestsError struct {
	// RetryAfter is how long the API asks us to wait before sending again;
	// zero if it doesn't say.
	RetryAfter time.Duration
}

func (e *TooManyRequestsError) Error() string {
	if e.RetryAfter > 0 {
		return fmt.Sprintf("429 Too Many Requests (retry after %v)", e.RetryAfter)
	}
	return "429 Too Many Requests"
}

// ThirdPartyAPI is the third-party API. Post sends a single message.
type ThirdPartyAPI interface {
	Post(ctx context.Context, message string) error
}

// MockAPI is a configurable fake of the third-party API.
//
// Every call waits for the configured delay, then returns the next scripted
// error (see SetErrors), or nil (success) once those run out. Messages that
// were posted successfully are recorded in order.
// The zero value is not usable; use NewMockAPI.
type MockAPI struct {
	mu     sync.Mutex
	delay  time.Duration
	errs   []error
	posted []string
}

// NewMockAPI creates a mock that responds immediately with success.
func NewMockAPI() *MockAPI {
	return &MockAPI{}
}

// SetDelay sets how long each Post call waits before responding.
func (m *MockAPI) SetDelay(d time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.delay = d
}

// SetErrors scripts what the next Post calls respond with, one per call, in
// order (nil means success). Once they run out, calls respond with nil.
func (m *MockAPI) SetErrors(errs ...error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.errs = append([]error(nil), errs...)
}

// Posted returns the messages that were posted successfully, in order.
func (m *MockAPI) Posted() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return append([]string(nil), m.posted...)
}

// Post waits for the configured delay (or until ctx is done), then responds
// with the next scripted error, or nil if there is none. A scripted error is
// used up even if ctx ends first.
func (m *MockAPI) Post(ctx context.Context, message string) error {
	m.mu.Lock()
	delay := m.delay
	var err error
	if len(m.errs) > 0 {
		err, m.errs = m.errs[0], m.errs[1:]
	}
	m.mu.Unlock()

	select {
	case <-time.After(delay):
	case <-ctx.Done():
		return ctx.Err()
	}

	if err == nil {
		m.mu.Lock()
		m.posted = append(m.posted, message)
		m.mu.Unlock()
	}
	return err
}
