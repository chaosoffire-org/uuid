package internal

import "testing"

// FuzzNewRandom tests that NewRandom works correctly under fuzzing.
func FuzzNewRandom(f *testing.F) {
	// Seed with some iteration counts
	f.Add(uint8(1))
	f.Add(uint8(10))
	f.Add(uint8(100))

	f.Fuzz(func(t *testing.T, iterations uint8) {
		// Generate multiple UUIDs and verify properties
		for i := uint8(0); i < iterations%50+1; i++ {
			uuid, err := NewRandom()
			if err != nil {
				t.Fatalf("NewRandom() error = %v", err)
			}

			// Verify version
			if uuid.Version() != V4 {
				t.Errorf("NewRandom() Version = %v, want V4", uuid.Version())
			}

			// Verify variant
			if uuid.Variant() != VariantRFC4122 {
				t.Errorf("NewRandom() Variant = %v, want VariantRFC4122", uuid.Variant())
			}

			// Verify roundtrip
			str := uuid.String()

			parsed, err := Parse(str)
			if err != nil {
				t.Errorf("Parse(uuid.String()) failed: %v", err)
			}

			if parsed != uuid {
				t.Errorf("roundtrip mismatch: %v != %v", uuid, parsed)
			}
		}
	})
}
