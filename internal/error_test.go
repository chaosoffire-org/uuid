package internal

import (
	"errors"
	"testing"
)

func TestInvalidLengthError_Error(t *testing.T) {
	tests := []struct {
		len  int
		want string
	}{
		{0, "invalid UUID length: 0"},
		{10, "invalid UUID length: 10"},
		{35, "invalid UUID length: 35"},
		{100, "invalid UUID length: 100"},
	}

	for _, tt := range tests {
		err := invalidLengthError{len: tt.len}
		if got := err.Error(); got != tt.want {
			t.Errorf("invalidLengthError{%d}.Error() = %q, want %q", tt.len, got, tt.want)
		}
	}
}

func TestInvalidLengthError_Is(t *testing.T) {
	err := invalidLengthError{len: 10}

	if !errors.Is(err, ErrInvalidLength) {
		t.Error("invalidLengthError should match ErrInvalidLength")
	}

	if errors.Is(err, ErrInvalidUUIDFormat) {
		t.Error("invalidLengthError should not match ErrInvalidUUIDFormat")
	}
}

func TestIsInvalidLengthError(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"invalidLengthError", invalidLengthError{len: 10}, true},
		{"ErrInvalidLength", ErrInvalidLength, true},
		{"wrapped invalidLengthError", errors.New("wrapped"), false},
		{"nil", nil, false},
		{"other error", ErrInvalidUUIDFormat, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidLengthError(tt.err); got != tt.want {
				t.Errorf("IsInvalidLengthError() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInvalidFormatError_Error(t *testing.T) {
	err := invalidFormatError{}
	want := "invalid UUID format"

	if got := err.Error(); got != want {
		t.Errorf("invalidFormatError.Error() = %q, want %q", got, want)
	}
}

func TestURNPrefixError_Error(t *testing.T) {
	tests := []struct {
		prefix string
		want   string
	}{
		{"", "invalid urn prefix"},
		{"urn:foo:", `invalid urn prefix: "urn:foo:"`},
		{"URN:UUID:", `invalid urn prefix: "URN:UUID:"`},
	}

	for _, tt := range tests {
		err := URNPrefixError{prefix: tt.prefix}
		if got := err.Error(); got != tt.want {
			t.Errorf("URNPrefixError{%q}.Error() = %q, want %q", tt.prefix, got, tt.want)
		}
	}
}

func TestURNPrefixError_Is(t *testing.T) {
	err := URNPrefixError{prefix: "test"}

	if !errors.Is(err, ErrInvalidURNPrefix) {
		t.Error("URNPrefixError should match ErrInvalidURNPrefix")
	}

	if errors.Is(err, ErrInvalidLength) {
		t.Error("URNPrefixError should not match ErrInvalidLength")
	}
}

func TestInvalidBracketedFormatError_Error(t *testing.T) {
	err := invalidBracketedFormatError{}
	want := "invalid bracketed UUID format"

	if got := err.Error(); got != want {
		t.Errorf("invalidBracketedFormatError.Error() = %q, want %q", got, want)
	}
}

func TestUnsupportedScanTypeError_Error(t *testing.T) {
	err := unsupportedScanTypeError{}
	want := "unsupported Scan type"

	if got := err.Error(); got != want {
		t.Errorf("unsupportedScanTypeError.Error() = %q, want %q", got, want)
	}
}

// Test sentinel errors are properly exported and usable
func TestSentinelErrors(t *testing.T) {
	tests := []struct {
		name string
		err  error
	}{
		{"ErrInvalidLength", ErrInvalidLength},
		{"ErrInvalidUUIDFormat", ErrInvalidUUIDFormat},
		{"ErrInvalidURNPrefix", ErrInvalidURNPrefix},
		{"ErrInvalidBracketedFormat", ErrInvalidBracketedFormat},
		{"ErrUnsupportedScanType", ErrUnsupportedScanType},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err == nil {
				t.Errorf("%s should not be nil", tt.name)
			}

			// Ensure Error() doesn't panic
			_ = tt.err.Error()
		})
	}
}
