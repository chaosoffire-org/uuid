package internal

import (
	"bytes"
	"errors"
	"testing"
)

func TestUUIDString(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")
	got := uuid.String()

	want := "12345678-1234-1234-1234-123456789012"
	if got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestUUIDURN(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")
	got := uuid.URN()

	want := "urn:uuid:12345678-1234-1234-1234-123456789012"
	if got != want {
		t.Errorf("URN() = %q, want %q", got, want)
	}
}

func TestUUIDVersion(t *testing.T) {
	tests := []struct {
		uuid    string
		version Version
	}{
		{"00000000-0000-1000-8000-000000000000", V1},
		{"00000000-0000-2000-8000-000000000000", V2},
		{"00000000-0000-3000-8000-000000000000", V3},
		{"00000000-0000-4000-8000-000000000000", V4},
		{"00000000-0000-5000-8000-000000000000", V5},
		{"00000000-0000-6000-8000-000000000000", V6},
		{"00000000-0000-7000-8000-000000000000", V7},
	}

	for _, tt := range tests {
		t.Run(tt.version.String(), func(t *testing.T) {
			uuid := MustParse(tt.uuid)
			if uuid.Version() != tt.version {
				t.Errorf("Version() = %v, want %v", uuid.Version(), tt.version)
			}
		})
	}
}

func TestUUIDVariant(t *testing.T) {
	tests := []struct {
		name    string
		uuid    UUID
		variant Variant
	}{
		{"RFC4122", MustParse("00000000-0000-4000-8000-000000000000"), VariantRFC4122},
		{"RFC4122-bf", MustParse("00000000-0000-4000-bf00-000000000000"), VariantRFC4122},
		{"NCS", UUID{0, 0, 0, 0, 0, 0, 0, 0, 0x00, 0, 0, 0, 0, 0, 0, 0}, VariantNCS},
		{"Microsoft", UUID{0, 0, 0, 0, 0, 0, 0, 0, 0xc0, 0, 0, 0, 0, 0, 0, 0}, VariantMicrosoft},
		{"Future", UUID{0, 0, 0, 0, 0, 0, 0, 0, 0xe0, 0, 0, 0, 0, 0, 0, 0}, VariantFuture},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.uuid.Variant() != tt.variant {
				t.Errorf("Variant() = %v, want %v", tt.uuid.Variant(), tt.variant)
			}
		})
	}
}

func TestSetVersion(t *testing.T) {
	versions := []Version{V1, V2, V3, V4, V5, V6, V7}
	for _, v := range versions {
		t.Run(v.String(), func(t *testing.T) {
			uuid := New()
			uuid.SetVersion(v)

			if uuid.Version() != v {
				t.Errorf("after SetVersion(%v), Version() = %v", v, uuid.Version())
			}
		})
	}
}

func TestSetVariant(t *testing.T) {
	variants := []Variant{VariantNCS, VariantRFC4122, VariantMicrosoft, VariantFuture}
	for _, v := range variants {
		t.Run(v.String(), func(t *testing.T) {
			uuid := New()
			uuid.SetVariant(v)

			if uuid.Variant() != v {
				t.Errorf("after SetVariant(%v), Variant() = %v", v, uuid.Variant())
			}
		})
	}
}

func TestBytes(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")

	b := uuid.Bytes()
	if len(b) != 16 {
		t.Errorf("Bytes() len = %d, want 16", len(b))
	}

	expected := []byte{0x12, 0x34, 0x56, 0x78, 0x12, 0x34, 0x12, 0x34, 0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12}
	if !bytes.Equal(b, expected) {
		t.Errorf("Bytes() = %x, want %x", b, expected)
	}
}

func TestMarshalText(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")

	text, err := uuid.MarshalText()
	if err != nil {
		t.Errorf("MarshalText() error = %v", err)
	}

	if string(text) != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("MarshalText() = %q, want %q", text, "12345678-1234-1234-1234-123456789012")
	}
}

func TestUnmarshalText(t *testing.T) {
	var uuid UUID

	err := uuid.UnmarshalText([]byte("12345678-1234-1234-1234-123456789012"))
	if err != nil {
		t.Errorf("UnmarshalText() error = %v", err)
	}

	expected := MustParse("12345678-1234-1234-1234-123456789012")
	if uuid != expected {
		t.Errorf("UnmarshalText() = %v, want %v", uuid, expected)
	}
}

func TestUnmarshalTextInvalid(t *testing.T) {
	var uuid UUID

	err := uuid.UnmarshalText([]byte("invalid"))
	if err == nil {
		t.Error("UnmarshalText(invalid) should return error")
	}
}

func TestMarshalBinary(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")

	data, err := uuid.MarshalBinary()
	if err != nil {
		t.Errorf("MarshalBinary() error = %v", err)
	}

	if len(data) != 16 {
		t.Errorf("MarshalBinary() len = %d, want 16", len(data))
	}
}

func TestUnmarshalBinary(t *testing.T) {
	original := MustParse("12345678-1234-1234-1234-123456789012")
	data, _ := original.MarshalBinary()

	var uuid UUID

	err := uuid.UnmarshalBinary(data)
	if err != nil {
		t.Errorf("UnmarshalBinary() error = %v", err)
	}

	if uuid != original {
		t.Errorf("UnmarshalBinary() = %v, want %v", uuid, original)
	}
}

func TestUnmarshalBinaryInvalidLength(t *testing.T) {
	var uuid UUID

	err := uuid.UnmarshalBinary([]byte{0, 1, 2})
	if !errors.Is(err, ErrInvalidLength) {
		t.Errorf("UnmarshalBinary(short) error = %v, want ErrInvalidLength", err)
	}
}

func TestValue(t *testing.T) {
	uuid := MustParse("12345678-1234-1234-1234-123456789012")

	val, err := uuid.Value()
	if err != nil {
		t.Errorf("Value() error = %v", err)
	}

	if val != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("Value() = %v, want 12345678-1234-1234-1234-123456789012", val)
	}
}

func TestScan(t *testing.T) {
	tests := []struct {
		name  string
		src   interface{}
		valid bool
	}{
		{"nil", nil, true},
		{"string", "12345678-1234-1234-1234-123456789012", true},
		{"bytes_36", []byte("12345678-1234-1234-1234-123456789012"), true},
		{"bytes_16", make([]byte, 16), true},
		{"invalid_string", "invalid", false},
		{"int", 123, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var uuid UUID

			err := uuid.Scan(tt.src)
			if (err == nil) != tt.valid {
				t.Errorf("Scan(%v) error = %v, valid = %v", tt.src, err, tt.valid)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		input   string
		valid   bool
		errType error
	}{
		{"12345678-1234-1234-1234-123456789012", true, nil},
		{"urn:uuid:12345678-1234-1234-1234-123456789012", true, nil},
		{"{12345678-1234-1234-1234-123456789012}", true, nil},
		{"12345678123412341234123456789012", true, nil},
		{"", false, ErrInvalidLength},
		{"invalid", false, ErrInvalidLength},
		{"12345678-1234-1234-1234-12345678901", false, ErrInvalidLength},
		{"urn:XXXX:12345678-1234-1234-1234-123456789012", false, ErrInvalidURNPrefix},
		{"12345678X1234-1234-1234-123456789012", false, ErrInvalidUUIDFormat},
		{"GGGGGGGG-GGGG-GGGG-GGGG-GGGGGGGGGGGG", false, ErrInvalidUUIDFormat},
		{"{12345678-1234-1234-1234-123456789012]", false, ErrInvalidBracketedFormat},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			_, err := Parse(tt.input)
			if tt.valid && err != nil {
				t.Errorf("Parse(%q) error = %v, want nil", tt.input, err)
			}

			if !tt.valid && err == nil {
				t.Errorf("Parse(%q) error = nil, want error", tt.input)
			}

			if tt.errType != nil && !errors.Is(err, tt.errType) {
				t.Errorf("Parse(%q) error = %v, want %v", tt.input, err, tt.errType)
			}
		})
	}
}

func TestParseBytes(t *testing.T) {
	input := []byte("12345678-1234-1234-1234-123456789012")

	uuid, err := ParseBytes(input)
	if err != nil {
		t.Errorf("ParseBytes() error = %v", err)
	}

	if uuid.String() != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("ParseBytes() = %v", uuid.String())
	}
}

func TestMustParse(t *testing.T) {
	// Valid
	uuid := MustParse("12345678-1234-1234-1234-123456789012")
	if uuid.String() != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("MustParse() = %v", uuid.String())
	}
}

func TestMustParsePanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("MustParse(invalid) should panic")
		}
	}()

	MustParse("invalid")
}

func TestFromBytes(t *testing.T) {
	data := []byte{0x12, 0x34, 0x56, 0x78, 0x12, 0x34, 0x12, 0x34, 0x12, 0x34, 0x12, 0x34, 0x56, 0x78, 0x90, 0x12}

	uuid, err := FromBytes(data)
	if err != nil {
		t.Errorf("FromBytes() error = %v", err)
	}

	if uuid.String() != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("FromBytes() = %v", uuid.String())
	}
}

func TestFromBytesInvalidLength(t *testing.T) {
	_, err := FromBytes([]byte{1, 2, 3})
	if !errors.Is(err, ErrInvalidLength) {
		t.Errorf("FromBytes(short) error = %v, want ErrInvalidLength", err)
	}
}

func TestMust(t *testing.T) {
	uuid := Must(Parse("12345678-1234-1234-1234-123456789012"))
	if uuid.String() != "12345678-1234-1234-1234-123456789012" {
		t.Errorf("Must() = %v", uuid.String())
	}
}

func TestMustPanic(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Must(error) should panic")
		}
	}()

	Must(Parse("invalid"))
}

func TestVersionString(t *testing.T) {
	tests := []struct {
		version Version
		want    string
	}{
		{V1, "VERSION_1"},
		{V2, "VERSION_2"},
		{V3, "VERSION_3"},
		{V4, "VERSION_4"},
		{V5, "VERSION_5"},
		{V6, "VERSION_6"},
		{V7, "VERSION_7"},
		{Version(8), "VERSION_8"},
		{Version(0), "VERSION_0"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.version.String(); got != tt.want {
				t.Errorf("Version.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVariantString(t *testing.T) {
	tests := []struct {
		variant Variant
		want    string
	}{
		{VariantNCS, "NCS"},
		{VariantRFC4122, "RFC4122"},
		{VariantMicrosoft, "Microsoft"},
		{VariantFuture, "Future"},
		{Variant(10), "Future"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			if got := tt.variant.String(); got != tt.want {
				t.Errorf("Variant.String() = %v, want %v", got, tt.want)
			}
		})
	}
}
