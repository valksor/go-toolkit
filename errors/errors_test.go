package errors

import (
	"errors"
	"fmt"
	"net/http"
	"testing"
)

func TestBaseError(t *testing.T) {
	err := NewBaseError(ErrorCodeNotFound, "resource not found")

	if err.Error() != "resource not found" {
		t.Errorf("BaseError.Error() = %v, want 'resource not found'", err.Error())
	}

	var baseErr *BaseError
	if !errors.As(err, &baseErr) {
		t.Error("NewBaseError() should return a *BaseError")
	}

	if baseErr.Code != ErrorCodeNotFound {
		t.Errorf("BaseError.Code = %v, want %v", baseErr.Code, ErrorCodeNotFound)
	}
}

func TestComponentError(t *testing.T) {
	baseErr := NewBaseError(ErrorCodeNoToken, "no token")
	compErr := NewComponentError("github", baseErr)

	if compErr == nil {
		t.Fatal("NewComponentError() returned nil")
	}

	expected := "github: no token"
	if compErr.Error() != expected {
		t.Errorf("ComponentError.Error() = %v, want %v", compErr.Error(), expected)
	}

	// Test Unwrap
	if !errors.Is(compErr, baseErr) {
		t.Error("ComponentError should unwrap to base error")
	}
}

func TestNewComponentError(t *testing.T) {
	// Test with nil error
	err := NewComponentError("test", nil)
	if err != nil {
		t.Errorf("NewComponentError() with nil error = %v, want nil", err)
	}

	// Test with actual error
	baseErr := NewBaseError(ErrorCodeInvalidConfig, "invalid config")
	err = NewComponentError("provider", baseErr)
	if err == nil {
		t.Fatal("NewComponentError() returned nil for non-nil error")
	}
}

func TestErrorTypeCheckers(t *testing.T) {
	tests := []struct {
		name      string
		err       error
		checkFunc func(error) bool
		expected  bool
	}{
		{"IsNoToken true", ErrNoToken, IsNoToken, true},
		{"IsNoToken false", ErrUnauthorized, IsNoToken, false},
		{"IsUnauthorized true", ErrUnauthorized, IsUnauthorized, true},
		{"IsUnauthorized false", ErrNotFound, IsUnauthorized, false},
		{"IsRateLimited true", ErrRateLimited, IsRateLimited, true},
		{"IsNotFound true", ErrNotFound, IsNotFound, true},
		{"IsInvalidReference true", ErrInvalidReference, IsInvalidReference, true},
		{"IsInvalidConfig true", ErrInvalidConfig, IsInvalidConfig, true},
		{"IsConflict true", ErrConflict, IsConflict, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if result := tt.checkFunc(tt.err); result != tt.expected {
				t.Errorf("%v(%v) = %v, want %v", tt.name, tt.err, result, tt.expected)
			}
		})
	}
}

func TestErrorTypeCheckersWithComponent(t *testing.T) {
	// Wrapped errors should still match
	err := NewComponentError("github", ErrNoToken)

	if !IsNoToken(err) {
		t.Error("IsNoToken() should return true for wrapped error")
	}

	if IsUnauthorized(err) {
		t.Error("IsUnauthorized() should return false for wrapped NoToken error")
	}
}

// Mock HTTP error with StatusCode.
type mockHTTPError struct {
	statusCode int
	msg        string
}

func (e *mockHTTPError) Error() string {
	return e.msg
}

func (e *mockHTTPError) StatusCode() int {
	return e.statusCode
}

func TestWrapHTTPError(t *testing.T) {
	tests := []struct {
		name          string
		err           error
		componentName string
		baseErrors    map[int]error
		checkFunc     func(error) bool
	}{
		{
			name:          "nil error",
			err:           nil,
			componentName: "test",
			baseErrors:    nil,
			checkFunc:     func(err error) bool { return err == nil },
		},
		{
			name:          "401 unauthorized",
			err:           &mockHTTPError{statusCode: http.StatusUnauthorized, msg: "unauthorized"},
			componentName: "github",
			baseErrors:    nil,
			checkFunc:     IsUnauthorized,
		},
		{
			name:          "404 not found",
			err:           &mockHTTPError{statusCode: http.StatusNotFound, msg: "not found"},
			componentName: "gitlab",
			baseErrors:    nil,
			checkFunc:     IsNotFound,
		},
		{
			name:          "409 conflict",
			err:           &mockHTTPError{statusCode: http.StatusConflict, msg: "conflict"},
			componentName: "jira",
			baseErrors:    nil,
			checkFunc:     IsConflict,
		},
		{
			name:          "429 rate limited",
			err:           &mockHTTPError{statusCode: http.StatusTooManyRequests, msg: "rate limited"},
			componentName: "api",
			baseErrors:    nil,
			checkFunc:     IsRateLimited,
		},
		{
			name:          "custom mapping",
			err:           &mockHTTPError{statusCode: http.StatusForbidden, msg: "forbidden"},
			componentName: "custom",
			baseErrors:    map[int]error{http.StatusForbidden: ErrInvalidConfig},
			checkFunc:     IsInvalidConfig,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := WrapHTTPError(tt.err, tt.componentName, tt.baseErrors)

			if !tt.checkFunc(result) {
				t.Errorf("WrapHTTPError() = %v, should match expected error type", result)
			}
		})
	}
}

func TestWrapHTTPErrorPreservesComponent(t *testing.T) {
	err := &mockHTTPError{statusCode: http.StatusUnauthorized, msg: "unauthorized"}
	wrapped := WrapHTTPError(err, "github", nil)

	var compErr *ComponentError
	if !errors.As(wrapped, &compErr) {
		t.Fatal("WrapHTTPError() should return a ComponentError")
	}

	if compErr.Component != "github" {
		t.Errorf("Component = %v, want 'github'", compErr.Component)
	}
}

func TestErrorCreationHelpers(t *testing.T) {
	tests := []struct {
		name      string
		fn        func() error
		checkFunc func(error) bool
	}{
		{"NoTokenError", func() error { return NoTokenError("github") }, IsNoToken},
		{"NotFoundError", func() error { return NotFoundError("gitlab", "project") }, IsNotFound},
		{"InvalidReferenceError", func() error { return InvalidReferenceError("jira", "PROJ-123") }, IsInvalidReference},
		{"InvalidConfigError", func() error { return InvalidConfigError("notion", "missing token") }, IsInvalidConfig},
		{"ConflictError", func() error { return ConflictError("azure", "update conflict") }, IsConflict},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.fn()

			if !tt.checkFunc(err) {
				t.Errorf("%s() should create matching error type", tt.name)
			}

			var compErr *ComponentError
			if !errors.As(err, &compErr) {
				t.Errorf("%s() should return a ComponentError", tt.name)
			}
		})
	}
}

func TestUnauthorizedError(t *testing.T) {
	baseErr := errors.New("authentication failed")
	err := UnauthorizedError("github", baseErr)

	if !IsUnauthorized(err) {
		t.Error("UnauthorizedError() should wrap ErrUnauthorized")
	}

	var compErr *ComponentError
	if !errors.As(err, &compErr) {
		t.Fatal("UnauthorizedError() should return a ComponentError")
	}

	if compErr.Component != "github" {
		t.Errorf("Component = %v, want 'github'", compErr.Component)
	}
}

func TestRateLimitedError(t *testing.T) {
	err := RateLimitedError("api", "too many requests")

	if !IsRateLimited(err) {
		t.Error("RateLimitedError() should wrap ErrRateLimited")
	}

	var compErr *ComponentError
	if !errors.As(err, &compErr) {
		t.Fatal("RateLimitedError() should return a ComponentError")
	}
}

func TestGetErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected ErrorCode
	}{
		{
			name:     "BaseError",
			err:      ErrNoToken,
			expected: ErrorCodeNoToken,
		},
		{
			name:     "wrapped BaseError",
			err:      NewComponentError("test", ErrUnauthorized),
			expected: ErrorCodeUnauthorized,
		},
		{
			name:     "generic error",
			err:      errors.New("generic error"),
			expected: ErrorCodeUnknown,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code := GetErrorCode(tt.err)
			if code != tt.expected {
				t.Errorf("GetErrorCode() = %v, want %v", code, tt.expected)
			}
		})
	}
}

func TestSentinelErrors(t *testing.T) {
	// Ensure all sentinel errors are properly created
	sentinelErrors := []struct {
		err       error
		errCode   ErrorCode
		checkFunc func(error) bool
	}{
		{ErrNoToken, ErrorCodeNoToken, IsNoToken},
		{ErrUnauthorized, ErrorCodeUnauthorized, IsUnauthorized},
		{ErrRateLimited, ErrorCodeRateLimited, IsRateLimited},
		{ErrNetworkError, ErrorCodeNetworkError, IsNetworkError},
		{ErrNotFound, ErrorCodeNotFound, IsNotFound},
		{ErrInvalidReference, ErrorCodeInvalidReference, IsInvalidReference},
		{ErrInvalidConfig, ErrorCodeInvalidConfig, IsInvalidConfig},
		{ErrConflict, ErrorCodeConflict, IsConflict},
	}

	for _, tt := range sentinelErrors {
		t.Run(fmt.Sprintf("%d", tt.errCode), func(t *testing.T) {
			if !tt.checkFunc(tt.err) {
				t.Errorf("Sentinel error %v should match its checker", tt.err)
			}

			if GetErrorCode(tt.err) != tt.errCode {
				t.Errorf("Sentinel error code = %v, want %v", GetErrorCode(tt.err), tt.errCode)
			}
		})
	}
}
