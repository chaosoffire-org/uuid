package internal

import (
	"errors"
	"testing"
	"unicode/utf8"
)

// FuzzParse tests that Parse handles any string input without panicking.
func FuzzParse(f *testing.F) {
	// Add seed corpus
	f.Add("c9436399-cc42-4e80-85ec-7bf4a8f76753")
	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("{6ba7b810-9dad-11d1-80b4-00c04fd430c8}")
	f.Add("urn:uuid:6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("6ba7b8109dad11d180b400c04fd430c8")
	f.Add("019c13Y044A077748Zc98zf445c77b65")
	f.Add("6BA7B810-9DAD-11D1-80B4-00C04FD430C8")

	f.Add("")
	f.Add("invalid")
	f.Add("12345678-1234-1234-1234-12345678901")   // too short (35)
	f.Add("12345678-1234-1234-1234-1234567890123") // too long (37)

	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430cZ")
	f.Add("6ba7b810-9dad-11d1-80b4-00c04fd430!8")
	f.Add("zzzzzzzz-zzzz-zzzz-zzzz-zzzzzzzzzzzz")
	f.Add("6ba7b810-XXXX-11d1-80b4-00c04fd430c8")

	f.Add("6ba7b81099dad-11d1-80b4-00c04fd430c8")
	f.Add("6ba7b810-9dad911d1-80b4-00c04fd430c8")
	f.Add("6ba7b810-9dad-11d180b4-00c04fd-430c8")

	f.Add("urn:xxxx:6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("{6ba7b810-9dad-11d1-80b4-00c04fd430c8")
	f.Add("{6ba7b810-9dad-11d1-80b4-00c04fd430c8z]")

	f.Add(string([]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15})) // binary data

	f.Fuzz(func(t *testing.T, input string) {
		// Should never panic
		uuid, err := Parse(input)
		if err == nil {
			// If parse succeeded, verify roundtrip
			str := uuid.String()

			parsed2, err2 := Parse(str)
			if err2 != nil {
				t.Errorf("roundtrip failed: Parse(%q).String() = %q, Parse again failed: %v", input, str, err2)
			}

			if parsed2 != uuid {
				t.Errorf("roundtrip mismatch: %v != %v", uuid, parsed2)
			}
		}
	})
}

// FuzzParseBytes tests that ParseBytes handles any byte slice input without panicking.
func FuzzParseBytes(f *testing.F) {
	// Add seed corpus
	f.Add([]byte("12345678-1234-1234-1234-123456789012"))
	f.Add([]byte{})
	f.Add([]byte("invalid"))
	f.Add(make([]byte, 16))
	f.Add(make([]byte, 36))
	f.Add([]byte{0xff, 0xff, 0xff, 0xff})

	f.Fuzz(func(t *testing.T, input []byte) {
		// Should never panic
		uuid, err := ParseBytes(input)
		if err == nil {
			// Valid parse - verify it produces a valid string
			str := uuid.String()
			if len(str) != 36 {
				t.Errorf("ParseBytes success but String() len = %d", len(str))
			}
		}
	})
}

// FuzzFromBytes tests that FromBytes handles any byte slice without panicking.
func FuzzFromBytes(f *testing.F) {
	// Add seed corpus
	f.Add(make([]byte, 16))
	f.Add([]byte{})
	f.Add(make([]byte, 15))
	f.Add(make([]byte, 17))
	f.Add([]byte{0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff, 0xff})

	f.Fuzz(func(t *testing.T, input []byte) {
		// Should never panic
		uuid, err := FromBytes(input)

		if len(input) == 16 {
			if err != nil {
				t.Errorf("FromBytes with 16 bytes should succeed, got error: %v", err)
			}
			// Verify roundtrip
			if string(uuid[:]) != string(input) {
				t.Errorf("FromBytes content mismatch")
			}
		} else {
			if !errors.Is(err, ErrInvalidLength) {
				t.Errorf("FromBytes with %d bytes should return ErrInvalidLength, got: %v", len(input), err)
			}
		}
	})
}

// FuzzUnmarshalText tests that UnmarshalText handles any byte slice without panicking.
func FuzzUnmarshalText(f *testing.F) {
	// Add seed corpus
	f.Add([]byte("12345678-1234-1234-1234-123456789012"))
	f.Add([]byte(""))
	f.Add([]byte("invalid"))

	f.Fuzz(func(t *testing.T, input []byte) {
		var uuid UUID
		// Should never panic
		err := uuid.UnmarshalText(input)
		if err == nil {
			// Verify the result can be marshaled back
			text, err2 := uuid.MarshalText()
			if err2 != nil {
				t.Errorf("MarshalText after successful UnmarshalText failed: %v", err2)
			}

			if !utf8.Valid(text) {
				t.Error("MarshalText produced invalid UTF-8")
			}
		}
	})
}

// FuzzUnmarshalBinary tests that UnmarshalBinary handles any byte slice without panicking.
func FuzzUnmarshalBinary(f *testing.F) {
	// Add seed corpus
	f.Add(make([]byte, 16))
	f.Add([]byte{})
	f.Add(make([]byte, 15))
	f.Add(make([]byte, 17))

	f.Fuzz(func(t *testing.T, input []byte) {
		var uuid UUID
		// Should never panic
		err := uuid.UnmarshalBinary(input)

		if len(input) == 16 {
			if err != nil {
				t.Errorf("UnmarshalBinary with 16 bytes should succeed, got error: %v", err)
			}
		} else {
			if err == nil {
				t.Errorf("UnmarshalBinary with %d bytes should fail", len(input))
			}
		}
	})
}

// FuzzScan tests that Scan handles various input types without panicking.
func FuzzScan(f *testing.F) {
	// Add seed corpus for strings
	f.Add("12345678-1234-1234-1234-123456789012")
	f.Add("")
	f.Add("invalid")

	f.Fuzz(func(t *testing.T, input string) {
		var uuid UUID
		// Should never panic with string input
		_ = uuid.Scan(input)
	})
}
