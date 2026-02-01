package internal

import "testing"

// FuzzNullUUIDScanString tests NullUUID.Scan with various string inputs.
func FuzzNullUUIDScanString(f *testing.F) {
	// Add seed corpus
	f.Add("12345678-1234-1234-1234-123456789012")
	f.Add("")
	f.Add("invalid")
	f.Add("urn:uuid:12345678-1234-1234-1234-123456789012")
	f.Add("{12345678-1234-1234-1234-123456789012}")
	f.Add("12345678123412341234123456789012")
	f.Add("GGGGGGGG-GGGG-GGGG-GGGG-GGGGGGGGGGGG")

	f.Fuzz(func(t *testing.T, input string) {
		var nu NullUUID
		// Should never panic
		err := nu.Scan(input)
		if err == nil {
			// If scan succeeded, UUID should be valid
			if !nu.Valid {
				t.Error("successful Scan should set Valid to true")
			}
		}
	})
}

// FuzzNullUUIDScanBytes tests NullUUID.Scan with various byte slice inputs.
func FuzzNullUUIDScanBytes(f *testing.F) {
	// Add seed corpus
	f.Add([]byte("12345678-1234-1234-1234-123456789012"))
	f.Add([]byte{})
	f.Add(make([]byte, 16))
	f.Add(make([]byte, 36))
	f.Add([]byte("invalid"))

	f.Fuzz(func(t *testing.T, input []byte) {
		var nu NullUUID
		// Should never panic
		_ = nu.Scan(input)
	})
}
