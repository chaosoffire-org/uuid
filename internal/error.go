package internal

import (
	"errors"
	"fmt"
)

// invalidLengthError is returned when a UUID has an invalid length.
type invalidLengthError struct{ len int }

func (err invalidLengthError) Error() string {
	return fmt.Sprintf("invalid UUID length: %d", err.len)
}

func (err invalidLengthError) Is(target error) bool {
	_, ok := target.(invalidLengthError)
	return ok
}

// ErrInvalidLength is the sentinel error for invalid length.
var ErrInvalidLength = invalidLengthError{}

// IsInvalidLengthError is matcher function for invalidLengthError.
func IsInvalidLengthError(err error) bool {
	return errors.Is(err, ErrInvalidLength)
}

// invalidFormatError is returned when a UUID has an invalid format.
type invalidFormatError struct{}

func (err invalidFormatError) Error() string {
	return "invalid UUID format"
}

// ErrInvalidUUIDFormat is the sentinel error for invalid format.
var ErrInvalidUUIDFormat = invalidFormatError{}

// URNPrefixError is returned when a UUID URN has an invalid prefix.
type URNPrefixError struct{ prefix string }

func (e URNPrefixError) Error() string {
	if e.prefix == "" {
		return "invalid urn prefix"
	}

	return fmt.Sprintf("invalid urn prefix: %q", e.prefix)
}

func (e URNPrefixError) Is(target error) bool {
	_, ok := target.(URNPrefixError)
	return ok
}

// ErrInvalidURNPrefix is the sentinel error for invalid URN prefix.
var ErrInvalidURNPrefix = URNPrefixError{}

// invalidBracketedFormatError is returned when a bracketed UUID has invalid format.
type invalidBracketedFormatError struct{}

func (err invalidBracketedFormatError) Error() string {
	return "invalid bracketed UUID format"
}

// ErrInvalidBracketedFormat is the sentinel error for invalid bracketed format.
var ErrInvalidBracketedFormat = invalidBracketedFormatError{}

// unsupportedScanTypeError is returned when Scan is called with an unsupported type.
type unsupportedScanTypeError struct{}

func (err unsupportedScanTypeError) Error() string {
	return "unsupported Scan type"
}

// ErrUnsupportedScanType is the sentinel error for unsupported scan type.
var ErrUnsupportedScanType = unsupportedScanTypeError{}
