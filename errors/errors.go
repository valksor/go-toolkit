// Package errors provides typed error types and utilities.
//
// This package provides a base error type system with error codes,
// component-specific error wrapping, and HTTP error handling.
//
// Basic usage:
//
//	err := errors.NoTokenError("github")
//	if errors.IsNoToken(err) {
//	    // Handle missing token
//	}
package errors

import (
	"errors"
	"fmt"
	"net"
	"net/http"
)

// Common error types that can be used across components.
// These are base errors that components can wrap with component-specific context.
var (
	// ErrNoToken is returned when no API token is configured.
	ErrNoToken = NewBaseError(ErrorCodeNoToken, "api token not found")

	// ErrUnauthorized is returned when the API token is invalid or expired.
	ErrUnauthorized = NewBaseError(ErrorCodeUnauthorized, "token unauthorized or expired")

	// ErrRateLimited is returned when the API rate limit is exceeded.
	ErrRateLimited = NewBaseError(ErrorCodeRateLimited, "api rate limit exceeded")

	// ErrNetworkError is returned for network-related errors.
	ErrNetworkError = NewBaseError(ErrorCodeNetworkError, "network error")

	// ErrNotFound is returned when a resource is not found.
	ErrNotFound = NewBaseError(ErrorCodeNotFound, "resource not found")

	// ErrInvalidReference is returned when a reference is invalid.
	ErrInvalidReference = NewBaseError(ErrorCodeInvalidReference, "invalid reference")

	// ErrInvalidConfig is returned when configuration is invalid.
	ErrInvalidConfig = NewBaseError(ErrorCodeInvalidConfig, "invalid configuration")

	// ErrConflict is returned when an update conflict occurs.
	ErrConflict = NewBaseError(ErrorCodeConflict, "update conflict")
)

// Error codes for categorizing errors.
type ErrorCode int

const (
	ErrorCodeUnknown ErrorCode = iota
	ErrorCodeNoToken
	ErrorCodeUnauthorized
	ErrorCodeRateLimited
	ErrorCodeNetworkError
	ErrorCodeNotFound
	ErrorCodeInvalidReference
	ErrorCodeInsufficientScope // For tokens with insufficient permissions
	ErrorCodeInvalidConfig     // For invalid configuration
	ErrorCodeConflict          // For update conflicts
)

// BaseError is a typed error that can be identified by code.
type BaseError struct {
	Msg  string
	Code ErrorCode
}

func (e *BaseError) Error() string {
	return e.Msg
}

// NewBaseError creates a new BaseError with the given code and message.
func NewBaseError(code ErrorCode, msg string) error {
	return &BaseError{Code: code, Msg: msg}
}

// ComponentError wraps an error with component name for better error messages.
type ComponentError struct {
	Err       error
	Component string
}

func (e *ComponentError) Error() string {
	return fmt.Sprintf("%s: %v", e.Component, e.Err)
}

func (e *ComponentError) Unwrap() error {
	return e.Err
}

// NewComponentError wraps an error with component context.
func NewComponentError(component string, err error) error {
	if err == nil {
		return nil
	}

	return &ComponentError{Component: component, Err: err}
}

// Type checking helpers

// IsNoToken returns true if err is or wraps ErrNoToken.
func IsNoToken(err error) bool {
	return errors.Is(err, ErrNoToken)
}

// IsUnauthorized returns true if err is or wraps ErrUnauthorized.
func IsUnauthorized(err error) bool {
	return errors.Is(err, ErrUnauthorized)
}

// IsRateLimited returns true if err is or wraps ErrRateLimited.
func IsRateLimited(err error) bool {
	return errors.Is(err, ErrRateLimited)
}

// IsNetworkError returns true if err is or wraps ErrNetworkError.
func IsNetworkError(err error) bool {
	return errors.Is(err, ErrNetworkError)
}

// IsNotFound returns true if err is or wraps ErrNotFound.
func IsNotFound(err error) bool {
	return errors.Is(err, ErrNotFound)
}

// IsInvalidReference returns true if err is or wraps ErrInvalidReference.
func IsInvalidReference(err error) bool {
	return errors.Is(err, ErrInvalidReference)
}

// IsInvalidConfig returns true if err is or wraps ErrInvalidConfig.
func IsInvalidConfig(err error) bool {
	return errors.Is(err, ErrInvalidConfig)
}

// IsConflict returns true if err is or wraps ErrConflict.
func IsConflict(err error) bool {
	return errors.Is(err, ErrConflict)
}

// WrapHTTPError converts HTTP status codes to typed errors.
// The componentName is used to create component-specific wrapped errors.
// The baseErrors map should contain status codes to base errors for component-specific mappings.
func WrapHTTPError(err error, componentName string, baseErrors map[int]error) error {
	if err == nil {
		return nil
	}

	// Check for network errors first
	var netErr net.Error
	if errors.As(err, &netErr) {
		return NewComponentError(componentName, fmt.Errorf("%w: %w", ErrNetworkError, err))
	}

	// Check for HTTP errors
	// Try common interfaces for HTTP status codes
	type statusCoder interface {
		StatusCode() int
	}
	type httpStatuser interface {
		HTTPStatusCode() int
	}

	var statusCode int
	if sc, ok := err.(statusCoder); ok {
		statusCode = sc.StatusCode()
	} else if hs, ok := err.(httpStatuser); ok {
		statusCode = hs.HTTPStatusCode()
	} else {
		// No status code available, return as-is
		return err
	}

	// Check component-specific mappings first
	if baseErr, ok := baseErrors[statusCode]; ok {
		return NewComponentError(componentName, fmt.Errorf("%w: %w", baseErr, err))
	}

	// Default mappings
	switch statusCode {
	case http.StatusUnauthorized:
		return NewComponentError(componentName, fmt.Errorf("%w: %w", ErrUnauthorized, err))
	case http.StatusForbidden:
		return NewComponentError(componentName, fmt.Errorf("%w: %w", ErrRateLimited, err))
	case http.StatusNotFound:
		return NewComponentError(componentName, fmt.Errorf("%w: %w", ErrNotFound, err))
	case http.StatusConflict:
		return NewComponentError(componentName, fmt.Errorf("%w: %w", ErrConflict, err))
	case http.StatusTooManyRequests:
		return NewComponentError(componentName, fmt.Errorf("%w: %w", ErrRateLimited, err))
	default:
		return err
	}
}

// Error creation helpers

// NoTokenError creates a component-specific "no token" error.
func NoTokenError(component string) error {
	return NewComponentError(component, ErrNoToken)
}

// UnauthorizedError creates a component-specific unauthorized error.
func UnauthorizedError(component string, err error) error {
	return NewComponentError(component, fmt.Errorf("%w: %w", ErrUnauthorized, err))
}

// RateLimitedError creates a component-specific rate limit error.
func RateLimitedError(component string, detail string) error {
	return NewComponentError(component, fmt.Errorf("%w: %s", ErrRateLimited, detail))
}

// NotFoundError creates a component-specific not found error.
func NotFoundError(component string, resource string) error {
	return NewComponentError(component, fmt.Errorf("%w: %s", ErrNotFound, resource))
}

// InvalidReferenceError creates a component-specific invalid reference error.
func InvalidReferenceError(component string, ref string) error {
	return NewComponentError(component, fmt.Errorf("%w: %s", ErrInvalidReference, ref))
}

// InvalidConfigError creates a component-specific invalid configuration error.
func InvalidConfigError(component string, detail string) error {
	return NewComponentError(component, fmt.Errorf("%w: %s", ErrInvalidConfig, detail))
}

// ConflictError creates a component-specific conflict error.
func ConflictError(component string, detail string) error {
	return NewComponentError(component, fmt.Errorf("%w: %s", ErrConflict, detail))
}

// GetErrorCode returns the error code if the error is a BaseError.
func GetErrorCode(err error) ErrorCode {
	var baseErr *BaseError
	if errors.As(err, &baseErr) {
		return baseErr.Code
	}

	return ErrorCodeUnknown
}
