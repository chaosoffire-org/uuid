// Package uuid provides RFC 4122 UUID implementation.
package uuid

import "github.com/chaosoffire-org/uuid-go/internal"

// Type aliases for public API
type (
	UUID     = internal.UUID
	NullUUID = internal.NullUUID
	UUIDs    = internal.UUIDs
	Version  = internal.Version
	Variant  = internal.Variant
)

// UUID versions
const (
	V1 = internal.V1
	V2 = internal.V2
	V3 = internal.V3
	V4 = internal.V4
	V5 = internal.V5
	V6 = internal.V6
	V7 = internal.V7
)

// UUID variants
const (
	VariantNCS       = internal.VariantNCS
	VariantRFC4122   = internal.VariantRFC4122
	VariantMicrosoft = internal.VariantMicrosoft
	VariantFuture    = internal.VariantFuture
)

// Nil is the nil UUID.
var Nil = internal.Nil

// URNPrefixError is the error type for invalid URN prefix
type URNPrefixError = internal.URNPrefixError

// Sentinel errors
var (
	ErrInvalidLength          = internal.ErrInvalidLength
	ErrInvalidUUIDFormat      = internal.ErrInvalidUUIDFormat
	ErrInvalidURNPrefix       = internal.ErrInvalidURNPrefix
	ErrInvalidBracketedFormat = internal.ErrInvalidBracketedFormat
	ErrUnsupportedScanType    = internal.ErrUnsupportedScanType
)

// IsInvalidLengthError reports whether err is an invalid length error.
func IsInvalidLengthError(err error) bool {
	return internal.IsInvalidLengthError(err)
}

// New creates a new random UUID v4.
func New() UUID {
	return internal.New()
}

// NewString creates a new random UUID v4 string.
func NewString() string {
	return internal.NewString()
}

// NewRandom creates a new random UUID v4, returning error on failure.
func NewRandom() (UUID, error) {
	return internal.NewRandom()
}

// NewV7 creates a new UUID v7 (Unix Epoch time-based with random bits).
func NewV7() (UUID, error) {
	return internal.NewV7()
}

// Parse decodes s into a UUID.
func Parse(s string) (UUID, error) {
	return internal.Parse(s)
}

// ParseBytes parses a byte slice into a UUID.
func ParseBytes(b []byte) (UUID, error) {
	return internal.ParseBytes(b)
}

// MustParse is like Parse but panics on error.
func MustParse(s string) UUID {
	return internal.MustParse(s)
}

// FromBytes creates a UUID from a byte slice.
func FromBytes(b []byte) (UUID, error) {
	return internal.FromBytes(b)
}

// Must panics if err is not nil.
func Must(uuid UUID, err error) UUID {
	return internal.Must(uuid, err)
}

// Validate returns nil if s is a valid UUID string.
func Validate(s string) error {
	_, err := internal.Parse(s)
	return err
}
