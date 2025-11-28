package errors

import (
	"errors"
	"fmt"
	"testing"
)

func TestCustomError_Error(t *testing.T) {
	tests := []struct {
		name     string
		err      *CustomError
		expected string
	}{
		{
			name:     "Simple Error",
			err:      New(100001, "Simple error"),
			expected: "[100001] Simple error",
		},
		{
			name:     "Wrapped Error",
			err:      Wrap(errors.New("inner error"), New(100001, "Wrapped error")),
			expected: "[100001] Wrapped error: inner error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.expected {
				t.Errorf("CustomError.Error() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestCustomError_Unwrap(t *testing.T) {
	inner := errors.New("inner error")
	err := Wrap(inner, ErrInternalServer)

	if !errors.Is(err, inner) {
		t.Errorf("errors.Is(err, inner) = false, want true")
	}

	unwrapped := errors.Unwrap(err)
	if unwrapped != inner {
		t.Errorf("errors.Unwrap(err) = %v, want %v", unwrapped, inner)
	}
}

func TestErrorCodes(t *testing.T) {
	// Verify some known error codes
	if ErrInternalServer.Code != 100001 {
		t.Errorf("ErrInternalServer code = %d, want 100001", ErrInternalServer.Code)
	}
	if ErrUnauthorized.Code != 200001 {
		t.Errorf("ErrUnauthorized code = %d, want 200001", ErrUnauthorized.Code)
	}
	if ErrForbidden.Code != 300001 {
		t.Errorf("ErrForbidden code = %d, want 300001", ErrForbidden.Code)
	}
	if ErrInvalidInput.Code != 400001 {
		t.Errorf("ErrInvalidInput code = %d, want 400001", ErrInvalidInput.Code)
	}
	if ErrDatabase.Code != 500001 {
		t.Errorf("ErrDatabase code = %d, want 500001", ErrDatabase.Code)
	}
	if ErrExternalAPI.Code != 500006 {
		t.Errorf("ErrExternalAPI code = %d, want 500006", ErrExternalAPI.Code)
	}
}

func TestWrap(t *testing.T) {
	baseErr := ErrInvalidInput
	innerErr := fmt.Errorf("field 'id' is missing")

	wrapped := Wrap(innerErr, baseErr)

	if wrapped.Code != baseErr.Code {
		t.Errorf("Wrapped error code = %d, want %d", wrapped.Code, baseErr.Code)
	}
	if wrapped.Message != baseErr.Message {
		t.Errorf("Wrapped error message = %s, want %s", wrapped.Message, baseErr.Message)
	}
	if wrapped.Err != innerErr {
		t.Errorf("Wrapped inner error = %v, want %v", wrapped.Err, innerErr)
	}
}
