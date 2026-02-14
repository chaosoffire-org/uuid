package uuid

import (
	"bytes"
	"testing"
)

func TestIsInvalidLengthError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{"invalid length error", ErrInvalidLength, true},
		{"invalid format error", ErrInvalidUUIDFormat, false},
		{"invalid urn prefix error", ErrInvalidURNPrefix, false},
		{"invalid bracketed format error", ErrInvalidBracketedFormat, false},
		{"unsupported scan type error", ErrUnsupportedScanType, false},
		{"nil error", nil, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsInvalidLengthError(tt.err); got != tt.expected {
				t.Errorf("IsInvalidLengthError() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestNew(t *testing.T) {
	u := New()
	if u.Version() != V4 {
		t.Errorf("New() returned UUID with version %v, want %v", u.Version(), V4)
	}
	if u.Variant() != VariantRFC4122 {
		t.Errorf("New() returned UUID with variant %v, want %v", u.Variant(), VariantRFC4122)
	}
	if u == Nil {
		t.Errorf("New() returned Nil UUID")
	}
}

func TestNewString(t *testing.T) {
	s := NewString()
	if len(s) != 36 {
		t.Errorf("NewString() returned length %d, want 36", len(s))
	}
	if _, err := Parse(s); err != nil {
		t.Errorf("NewString() returned invalid UUID string: %v", err)
	}
}

func TestNewRandom(t *testing.T) {
	u, err := NewRandom()
	if err != nil {
		t.Fatalf("NewRandom() failed: %v", err)
	}
	if u.Version() != V4 {
		t.Errorf("NewRandom() returned UUID with version %v, want %v", u.Version(), V4)
	}
	if u.Variant() != VariantRFC4122 {
		t.Errorf("NewRandom() returned UUID with variant %v, want %v", u.Variant(), VariantRFC4122)
	}
}

func TestNewV7(t *testing.T) {
	u, err := NewV7()
	if err != nil {
		t.Fatalf("NewV7() failed: %v", err)
	}
	if u.Version() != V7 {
		t.Errorf("NewV7() returned UUID with version %v, want %v", u.Version(), V7)
	}
	if u.Variant() != VariantRFC4122 {
		t.Errorf("NewV7() returned UUID with variant %v, want %v", u.Variant(), VariantRFC4122)
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid UUID", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", false},
		{"valid UUID with braces", "{6ba7b810-9dad-11d1-80b4-00c04fd430c8}", false},
		{"valid UUID with urn prefix", "urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8", false},
		{"invalid length", "6ba7b810-9dad-11d1-80b4-00c04fd430c", true},
		{"invalid char", "6ba7b810-9dad-11d1-80b4-00c04fd430cz", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Parse(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestParseBytes(t *testing.T) {
	validUUID := MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{"valid UUID bytes", []byte(validUUID.String()), false},
		{"invalid length", validUUID[:15], true},
		{"nil bytes", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := ParseBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMustParse(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("MustParse() did not panic on invalid input")
		}
	}()
	MustParse("invalid-uuid")
}

func TestFromBytes(t *testing.T) {
	validUUID := MustParse("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	tests := []struct {
		name    string
		input   []byte
		wantErr bool
	}{
		{"valid bytes", validUUID[:], false},
		{"invalid length", []byte{1, 2, 3}, true},
		{"nil bytes", nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u, err := FromBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromBytes() error = %v, wantErr %v", err, tt.wantErr)
			}
			if err == nil && !bytes.Equal(u[:], tt.input) {
				t.Errorf("FromBytes() returned mismatching bytes")
			}
		})
	}
}

func TestMust(t *testing.T) {
	t.Run("no error", func(t *testing.T) {
		u := Must(New(), nil)
		if u == Nil {
			t.Errorf("Must() returned Nil on success")
		}
	})

	t.Run("with error", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Must() did not panic on error")
			}
		}()
		Must(Nil, ErrInvalidLength)
	})
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{"valid", "6ba7b810-9dad-11d1-80b4-00c04fd430c8", false},
		{"invalid", "invalid-uuid", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := Validate(tt.input); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
