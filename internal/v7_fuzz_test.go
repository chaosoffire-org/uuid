package internal

import "testing"

// FuzzNewV7 tests that NewV7 works correctly under fuzzing.
func FuzzNewV7(f *testing.F) {
	f.Add(uint8(1))
	f.Add(uint8(10))
	f.Add(uint8(100))

	f.Fuzz(func(t *testing.T, iterations uint8) {
		for i := uint8(0); i < iterations%50+1; i++ {
			uuid, err := NewV7()
			if err != nil {
				t.Fatalf("NewV7() error = %v", err)
			}

			if uuid.Version() != V7 {
				t.Errorf("NewV7() Version = %v, want V7", uuid.Version())
			}

			if uuid.Variant() != VariantRFC4122 {
				t.Errorf("NewV7() Variant = %v, want VariantRFC4122", uuid.Variant())
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
